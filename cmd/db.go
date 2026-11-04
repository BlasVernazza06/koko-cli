package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/BlasVernazza06/koko-cli/internal/features/db"
	"github.com/spf13/cobra"
)

var (
	dbDirFlag    string
	dbDryRunFlag bool
)

// dbCmd representa el comando raíz `koko db`.
var dbCmd = &cobra.Command{
	Use:   "db [command]",
	Short: "Unified database operations wrapper (migrate, push, generate, studio, seed, reset)",
	Long: `koko db is an intelligent command wrapper for database operations.
It inspects koko.config.json or package.json, automatically detects your ORM (Drizzle,
Prisma, Mongoose), layout (monorepo or standalone), and package manager (pnpm, npm, bun),
and runs the appropriate command with full terminal interactivity.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		targetDir := dbDirFlag
		if targetDir == "" {
			targetDir = "."
		}
		return db.PrintDBStackOverview(cmd.OutOrStdout(), targetDir)
	},
}

func makeDBSubcommand(action db.Action, shortDesc string) *cobra.Command {
	return &cobra.Command{
		Use:                string(action) + " [flags] [-- extra args]",
		Short:              shortDesc,
		DisableFlagParsing: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			targetDir := dbDirFlag
			if targetDir == "" {
				targetDir = "."
			}

			runner := db.NewRunner()
			opts := db.Options{
				Dir:    targetDir,
				Action: action,
				Args:   args,
				DryRun: dbDryRunFlag,
				Stdin:  os.Stdin,
				Stdout: cmd.OutOrStdout(),
				Stderr: cmd.ErrOrStderr(),
			}

			if err := runner.Run(context.Background(), opts); err != nil {
				return fmt.Errorf("koko db %s failed: %w", action, err)
			}
			return nil
		},
	}
}

func init() {
	dbCmd.PersistentFlags().StringVarP(&dbDirFlag, "dir", "d", ".", "Target project directory")
	dbCmd.PersistentFlags().BoolVar(&dbDryRunFlag, "dry-run", false, "Print the resolved command without executing it")

	// Subcomandos
	dbCmd.AddCommand(makeDBSubcommand(db.ActionMigrate, "Run pending database schema migrations"))
	dbCmd.AddCommand(makeDBSubcommand(db.ActionPush, "Push schema state directly to database (drizzle-kit push / prisma db push)"))
	dbCmd.AddCommand(makeDBSubcommand(db.ActionGenerate, "Generate ORM client and type definitions (drizzle-kit generate / prisma generate)"))
	dbCmd.AddCommand(makeDBSubcommand(db.ActionStudio, "Open interactive database visual UI / studio"))
	dbCmd.AddCommand(makeDBSubcommand(db.ActionSeed, "Seed database with test fixtures"))
	dbCmd.AddCommand(makeDBSubcommand(db.ActionReset, "Reset database schema and drop tables"))

	rootCmd.AddCommand(dbCmd)
}
