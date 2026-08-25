package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	kokoConfig "github.com/BlasVernazza06/koko-cli/internal/config"
	"github.com/BlasVernazza06/koko-cli/internal/scaffold"
	"github.com/BlasVernazza06/koko-cli/internal/validator"
	"github.com/spf13/cobra"
)

var (
	defaultFlag  bool
	frontendFlag string
	backendFlag  string
	pmFlag       string
	databaseFlag string
	ormFlag      string
	authFlag     string
	gitFlag      string
	recipieFlag  string
)

var initCmd = &cobra.Command{
	Use:   "init [name]",
	Short: "Initialize a new Koko project",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := ""
		if len(args) > 0 && args[0] != "" {
			projectName = args[0]
		}

		if defaultFlag {
			if projectName == "" {
				projectName = "my-project"
			}
			if err := validator.Validate(projectName); err != nil {
				fmt.Printf("\n\033[31m✗ Error: %s\033[0m\n\n", err.Error())
				os.Exit(1)
			}
			runDefaultInit(projectName, recipieFlag)
			return
		}
		if frontendFlag != "" || backendFlag != "" || databaseFlag != "" || ormFlag != "" || authFlag != "" {
			if projectName == "" {
				projectName = "my-project"
			}
			if err := validator.Validate(projectName); err != nil {
				fmt.Printf("\n\033[31m✗ Error: %s\033[0m\n\n", err.Error())
				os.Exit(1)
			}
			pm := pmFlag
			if pm == "" {
				pm = "pnpm"
			}
			initGit := strings.ToLower(gitFlag) == "yes" || strings.ToLower(gitFlag) == "true" || strings.ToLower(gitFlag) == "y"
			cfg := scaffold.ScaffoldConfig{
				ProjectName:    projectName,
				Frontend:       frontendFlag,
				Backend:        backendFlag,
				PackageManager: pm,
				Database:       databaseFlag,
				ORM:            ormFlag,
				Auth:           authFlag,
				InitGit:        initGit,
			}
			runManualInit(cfg)
			return
		}

		RunTUIInit(projectName)
	},
}

func runDefaultInit(projectName string, recipie string) {
	fmt.Printf("\n\033[90m┌\033[0m  \033[1mCreating a new Koko project\033[0m\n")
	fmt.Printf("\033[90m│\033[0m\n")
	fmt.Printf("\033[90m│\033[0m  \033[38;2;0;255;127m◆\033[0m  \033[1mProject name\033[0m     \033[90m·\033[0m  \033[38;2;167;139;250m%s\033[0m\n", projectName)
	fmt.Printf("\033[90m│\033[0m  \033[38;2;0;255;127m◆\033[0m  \033[1mRecipe\033[0m           \033[90m·\033[0m  \033[38;2;167;139;250mSaaS Starter\033[0m\n")
	fmt.Printf("\033[90m│\033[0m  \033[38;2;0;255;127m◆\033[0m  \033[1mPackage Manager\033[0m  \033[90m·\033[0m  \033[38;2;167;139;250mpnpm\033[0m\n")
	fmt.Printf("\033[90m│\033[0m\n")

	startTime := time.Now()

	cfg := scaffold.ScaffoldConfig{
		ProjectName: projectName,
		Recipe:      recipie,
		InitGit:     true,
	}

	steps := []struct {
		name string
		fn   func() error
	}{
		{"Scaffolding directories and copying templates...", func() error { return scaffold.CopyTemplates(projectName, cfg) }},
		{"Generating Docker and Database configuration...", func() error { return scaffold.GenerateDockerAndDB(projectName, cfg) }},
		{"Initializing Git repository...", func() error { return scaffold.InitGit(projectName) }},
		{"Creating koko.config.json manifest...", func() error { return kokoConfig.GenerateConfig(projectName, cfg) }},
	}

	for _, step := range steps {
		if err := step.fn(); err != nil {
			fmt.Printf("\033[90m│\033[0m  \033[31m✗\033[0m  \033[31m%s\033[0m\n", step.name)
			fmt.Printf("\033[90m│\033[0m\n")
			fmt.Printf("\033[90m└\033[0m  \033[31m\033[1mCreation error: %v\033[0m\n\n", err)
			os.Exit(1)
		}
		fmt.Printf("\033[90m│\033[0m  \033[32m✓\033[0m  %s\n", step.name)
	}

	elapsed := time.Since(startTime)
	fmt.Printf("\033[90m│\033[0m\n")
	fmt.Printf("\033[90m└\033[0m  \033[32m\033[1mProject created successfully in %.2fs!\033[0m\n\n", elapsed.Seconds())
	fmt.Println("  \033[1mNext steps:\033[0m")
	fmt.Printf("  1. \033[38;2;90;79;196mcd %s\033[0m\n", projectName)
	fmt.Println("  2. \033[38;2;90;79;196mpnpm install\033[0m")
	fmt.Println("  3. \033[38;2;90;79;196mpnpm dev\033[0m")
	fmt.Println()
}

func runManualInit(cfg scaffold.ScaffoldConfig) {

	fmt.Printf("\n\033[90m┌\033[0m  \033[1mCreating a new Koko project\033[0m\n")
	fmt.Printf("\033[90m│\033[0m\n")
	fmt.Printf("\033[90m│\033[0m  \033[38;2;0;255;127m◆\033[0m  \033[1mProject name\033[0m     \033[90m·\033[0m  \033[38;2;167;139;250m%s\033[0m\n", cfg.ProjectName)
	if cfg.Frontend != "" {
		fmt.Printf("\033[90m│\033[0m  \033[38;2;0;255;127m◆\033[0m  \033[1mFrontend\033[0m         \033[90m·\033[0m  \033[38;2;167;139;250m%s\033[0m\n", cfg.Frontend)
	}
	if cfg.Backend != "" {
		fmt.Printf("\033[90m│\033[0m  \033[38;2;0;255;127m◆\033[0m  \033[1mBackend\033[0m          \033[90m·\033[0m  \033[38;2;167;139;250m%s\033[0m\n", cfg.Backend)
	}
	if cfg.Database != "" {
		fmt.Printf("\033[90m│\033[0m  \033[38;2;0;255;127m◆\033[0m  \033[1mDatabase\033[0m         \033[90m·\033[0m  \033[38;2;167;139;250m%s\033[0m\n", cfg.Database)
	}
	if cfg.ORM != "" {
		fmt.Printf("\033[90m│\033[0m  \033[38;2;0;255;127m◆\033[0m  \033[1mORM\033[0m              \033[90m·\033[0m  \033[38;2;167;139;250m%s\033[0m\n", cfg.ORM)
	}
	if cfg.Auth != "" {
		fmt.Printf("\033[90m│\033[0m  \033[38;2;0;255;127m◆\033[0m  \033[1mAuth\033[0m             \033[90m·\033[0m  \033[38;2;167;139;250m%s\033[0m\n", cfg.Auth)
	}
	fmt.Printf("\033[90m│\033[0m  \033[38;2;0;255;127m◆\033[0m  \033[1mPackage Manager\033[0m  \033[90m·\033[0m  \033[38;2;167;139;250m%s\033[0m\n", cfg.PackageManager)
	fmt.Printf("\033[90m│\033[0m\n")

	startTime := time.Now()

	steps := []struct {
		name string
		fn   func() error
	}{
		{"Scaffolding directories and copying templates...", func() error { return scaffold.CopyTemplates(cfg.ProjectName, cfg) }},
		{"Generating Docker and Database configuration...", func() error { return scaffold.GenerateDockerAndDB(cfg.ProjectName, cfg) }},
	}
	if cfg.InitGit {
		steps = append(steps, struct {
			name string
			fn   func() error
		}{"Initializing Git repository...", func() error { return scaffold.InitGit(cfg.ProjectName) }})
	}
	steps = append(steps, struct {
		name string
		fn   func() error
	}{"Creating koko.config.json manifest...", func() error { return kokoConfig.GenerateConfig(cfg.ProjectName, cfg) }})
	for _, step := range steps {
		if err := step.fn(); err != nil {
			fmt.Printf("\033[90m│\033[0m  \033[31m✗\033[0m  \033[31m%s\033[0m\n", step.name)
			fmt.Printf("\033[90m│\033[0m\n")
			fmt.Printf("\033[90m└\033[0m  \033[31m\033[1mCreation error: %v\033[0m\n\n", err)
			os.Exit(1)
		}
		fmt.Printf("\033[90m│\033[0m  \033[32m✓\033[0m  %s\n", step.name)
	}

	elapsed := time.Since(startTime)

	fmt.Printf("\033[90m│\033[0m\n")
	fmt.Printf("\033[90m└\033[0m  \033[32m\033[1mProject created successfully in %.2fs!\033[0m\n\n", elapsed.Seconds())
	fmt.Println("  \033[1mNext steps:\033[0m")
	fmt.Printf("  1. \033[38;2;90;79;196mcd %s\033[0m\n", cfg.ProjectName)
	fmt.Printf("  2. \033[38;2;90;79;196m%s install\033[0m\n", cfg.PackageManager)
	fmt.Printf("  3. \033[38;2;90;79;196m%s dev\033[0m\n", cfg.PackageManager)
	fmt.Println()
}

func init() {
	initCmd.Flags().BoolVarP(&defaultFlag, "default", "d", false, "Create project directly using default configuration and recipe (SaaS Starter)")

	initCmd.Flags().StringVarP(&frontendFlag, "frontend", "f", "", "Frontend Framework (ej: Nextjs, react, vue)")
	initCmd.Flags().StringVarP(&backendFlag, "backend", "b", "", "Backend Framework (ej: express, hono, fiber)")
	initCmd.Flags().StringVarP(&pmFlag, "package-manager", "p", "pnpm", "Package Manager to use (ej: bun, pnpm, npm, yarn)")
	initCmd.Flags().StringVar(&databaseFlag, "database", "", "Database (ej: postgres, mysql, mongodb)")
	initCmd.Flags().StringVar(&ormFlag, "orm", "", "ORM (ej: drizzle, prisma, mongoose)")
	initCmd.Flags().StringVar(&authFlag, "auth", "", "Auth Provider (ej: Better-auth, clerk, NextAuth.js)")
	initCmd.Flags().StringVar(&gitFlag, "git", "no", "Initialize Git Repository")

	initCmd.Flags().StringVarP(&recipieFlag, "recipie", "r", "saas", "Choose a recipie template (ej: SaaS, PERN, MERN, FAST_API_REACT)")

	rootCmd.AddCommand(initCmd)
}
