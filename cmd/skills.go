package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/BlasVernazza06/koko-cli/internal/features/skills"
	"github.com/spf13/cobra"
)

var (
	skillsDirFlag    string
	skillsYesFlag    bool
	skillsDryRunFlag bool
	skillsAgentFlag  string

	skillsManager     = skills.NewManager()
	skillsOrchestrator *skills.Orchestrator
)

var skillsCmd = &cobra.Command{
	Use:   "skills [args...]",
	Short: "Install and manage AI agent skills natively with autoskills integration",
	Long: `Install and manage curated AI agent skills (Cursor, Claude Code, GitHub Copilot)
natively with hybrid supply (GitHub raw, local user cache and embedded offline fallbacks).

Scans your project's configuration files (package.json, tsconfig, astro, wrangler) and physical file tree
to suggest and inject rules (.cursor/rules/*.mdc) and synchronize koko.config.json.

Examples:
  koko skills                  # Interactive Bubble Tea TUI
  koko skills --yes            # Auto-install recommended skills non-interactively
  koko skills --dry-run        # Preview detected/recommended skills without disk writes
  koko skills --dir ./my-app   # Target specific project directory
  koko skills ui               # Delegate to autoskills UI dashboard if arguments provided`,
	FParseErrWhitelist: cobra.FParseErrWhitelist{
		UnknownFlags: true,
	},
	Args:         cobra.ArbitraryArgs,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		targetDir := skillsDirFlag
		if targetDir == "" {
			targetDir = "."
		}

		ctx := cmd.Context()
		if ctx == nil {
			ctx = context.Background()
		}

		opts := skills.Options{
			Dir:    targetDir,
			Yes:    skillsYesFlag,
			DryRun: skillsDryRunFlag,
			Agent:  skillsAgentFlag,
			Args:   args,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
			Stdin:  os.Stdin,
		}

		// Si el usuario pasa subcomandos explícitos (como "ui" u opciones directas de autoskills)
		// y no es la invocación por defecto, delegar en autoskills runner si está disponible.
		if len(args) > 0 {
			runner, err := skillsManager.DetectRunner()
			if err == nil && runner != nil {
				return skillsManager.Run(ctx, opts)
			}
		}

		// Flujo Nativo de Koko Skills (Cero dependencia obligatoria de Node.js en runtime)
		orch := skillsOrchestrator
		if orch == nil {
			var err error
			orch, err = skills.NewOrchestrator()
			if err != nil {
				return fmt.Errorf("failed to initialize skills orchestrator: %w", err)
			}
		}

		return orch.ExecuteNativeSkills(ctx, opts)
	},
}

func init() {
	skillsCmd.Flags().StringVarP(&skillsDirFlag, "dir", "d", ".", "Target project directory to manage skills")
	skillsCmd.Flags().BoolVarP(&skillsYesFlag, "yes", "y", false, "Skip confirmation prompts (auto-accept recommended skills)")
	skillsCmd.Flags().BoolVar(&skillsDryRunFlag, "dry-run", false, "Show skills that would be installed without making changes")
	skillsCmd.Flags().StringVarP(&skillsAgentFlag, "agent", "a", "", "Target AI agent (e.g. cursor, claude, copilot)")

	rootCmd.AddCommand(skillsCmd)
}
