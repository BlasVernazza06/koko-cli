package cmd

import (
	"fmt"
	"os"
	"time"

	kokoConfig "github.com/BlasVernazza06/koko-cli/internal/config"
	"github.com/BlasVernazza06/koko-cli/internal/scaffold"
	"github.com/spf13/cobra"
)

var defaultFlag bool

var initCmd = &cobra.Command{
	Use:   "init [nombre]",
	Short: "Inicializa un nuevo proyecto Koko",
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
			runDefaultInit(projectName)
			return
		}

		RunTUIInit(projectName)
	},
}

func runDefaultInit(projectName string) {
	fmt.Printf("\n\033[90m┌\033[0m  \033[1mCreating a new Koko project\033[0m\n")
	fmt.Printf("\033[90m│\033[0m\n")
	fmt.Printf("\033[90m│\033[0m  \033[38;2;0;255;127m◆\033[0m  \033[1mProject name\033[0m     \033[90m·\033[0m  \033[38;2;167;139;250m%s\033[0m\n", projectName)
	fmt.Printf("\033[90m│\033[0m  \033[38;2;0;255;127m◆\033[0m  \033[1mRecipe\033[0m           \033[90m·\033[0m  \033[38;2;167;139;250mSaaS Starter\033[0m\n")
	fmt.Printf("\033[90m│\033[0m  \033[38;2;0;255;127m◆\033[0m  \033[1mPackage Manager\033[0m  \033[90m·\033[0m  \033[38;2;167;139;250mpnpm\033[0m\n")
	fmt.Printf("\033[90m│\033[0m\n")

	startTime := time.Now()

	cfg := scaffold.ScaffoldConfig{
		ProjectName: projectName,
		Recipe:      "saas",
	}

	steps := []struct {
		name string
		fn   func() error
	}{
		{"Estructurando directorios y copiando plantillas...", func() error { return scaffold.CopyTemplates(projectName, cfg) }},
		{"Generando configuración de Docker y DB...", func() error { return scaffold.GenerateDockerAndDB(projectName, cfg) }},
		{"Inicializando repositorio Git...", func() error { return scaffold.InitGit(projectName) }},
		{"Creando manifiesto koko.config.json...", func() error { return kokoConfig.GenerateConfig(projectName, cfg) }},
	}

	for _, step := range steps {
		if err := step.fn(); err != nil {
			fmt.Printf("\033[90m│\033[0m  \033[31m✗\033[0m  \033[31m%s\033[0m\n", step.name)
			fmt.Printf("\033[90m│\033[0m\n")
			fmt.Printf("\033[90m└\033[0m  \033[31m\033[1mError en la creación: %v\033[0m\n\n", err)
			os.Exit(1)
		}
		fmt.Printf("\033[90m│\033[0m  \033[32m✓\033[0m  %s\n", step.name)
	}

	elapsed := time.Since(startTime)
	fmt.Printf("\033[90m│\033[0m\n")
	fmt.Printf("\033[90m└\033[0m  \033[32m\033[1m¡Proyecto creado con éxito en %.2fs!\033[0m\n\n", elapsed.Seconds())
	fmt.Println("  \033[1mPróximos pasos:\033[0m")
	fmt.Printf("  1. \033[38;2;90;79;196mcd %s\033[0m\n", projectName)
	fmt.Println("  2. \033[38;2;90;79;196mpnpm install\033[0m")
	fmt.Println("  3. \033[38;2;90;79;196mpnpm dev\033[0m")
	fmt.Println()
}

func init() {
	initCmd.Flags().BoolVarP(&defaultFlag, "default", "d", false, "Crea el proyecto directamente con la configuración y receta por defecto (SaaS Starter)")
	rootCmd.AddCommand(initCmd)
}
