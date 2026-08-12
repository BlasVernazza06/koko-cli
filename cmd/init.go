package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	kokoConfig "github.com/BlasVernazza06/koko-cli/internal/config"
	"github.com/BlasVernazza06/koko-cli/internal/scaffold"
	"github.com/BlasVernazza06/koko-cli/internal/ui"
)

var autoApprove bool

func runInitSteps(projectName string, config scaffold.ScaffoldConfig) {
	// Step 1: Estructurando directorios y copiando plantillas
	s1 := ui.NewSpinner("[1/4] Estructurando directorios y copiando plantillas...")
	s1.Start()
	err := scaffold.CopyTemplates(projectName, config)
	time.Sleep(400 * time.Millisecond) // Premium feel/micro-animation delay
	if err != nil {
		s1.Stop(false)
		fmt.Printf("❌ Error al estructurar directorios: %v\n", err)
		os.Exit(1)
	}
	s1.Stop(true)

	// Step 2: Generando configuración de Docker y DB
	s2 := ui.NewSpinner("[2/4] Generando configuración de Docker y DB...")
	s2.Start()
	err = scaffold.GenerateDockerAndDB(projectName, config)
	time.Sleep(400 * time.Millisecond)
	if err != nil {
		s2.Stop(false)
		fmt.Printf("❌ Error al generar configuración de Docker y DB: %v\n", err)
		os.Exit(1)
	}
	s2.Stop(true)

	// Step 3: Inicializando repositorio Git
	s3 := ui.NewSpinner("[3/4] Inicializando repositorio Git...")
	s3.Start()
	err = scaffold.InitGit(projectName)
	time.Sleep(400 * time.Millisecond)
	if err != nil {
		s3.Stop(false)
		fmt.Printf("❌ Error al inicializar repositorio Git: %v\n", err)
		os.Exit(1)
	}
	s3.Stop(true)

	// Step 4: Creando manifiesto koko.config.json
	s4 := ui.NewSpinner("[4/4] Creando manifiesto koko.config.json...")
	s4.Start()
	err = kokoConfig.GenerateConfig(projectName, config)
	time.Sleep(400 * time.Millisecond)
	if err != nil {
		s4.Stop(false)
		fmt.Printf("❌ Error al crear manifiesto: %v\n", err)
		os.Exit(1)
	}
	s4.Stop(true)
}

var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize a new software project starter",
	Long:  `Run an interactive wizard to configure and generate a software project template.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := "my-project"
		if len(args) > 0 {
			projectName = args[0]
		}

		fmt.Println(`
			 __ _    __ _ __      __
			/  _ | |/ _' |\ \ /\ / /
			| (__| | (_| | \ V  V / 
			\___ |_|\__,_|  \_/\_/  v0.1.0
		`)

		// SIEMPRE preguntamos por el nombre al principio en modo interactivo
		if !autoApprove {
			nameForm := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("¿Cuál es el nombre de tu proyecto?").
						Value(&projectName),
				),
			)
			if err := nameForm.Run(); err != nil {
				fmt.Println("Error ejecutando formulario:", err)
				os.Exit(1)
			}
		}

		projectName = strings.TrimSpace(projectName)

		// Flujo 1: Inicialización automática con configuración por defecto (SaaS Starter)
		if autoApprove {
			fmt.Printf("🚀 Inicializando '%s' con configuración por defecto (Next.js + Drizzle + Better-Auth + Stripe)...\n", projectName)
			config := scaffold.ScaffoldConfig{
				ProjectName: projectName,
				Recipe:      "saas",
			}
			runInitSteps(projectName, config)
			fmt.Printf("\n✓ ¡Proyecto '%s' inicializado con éxito!\n", projectName)
			return
		}

		// Flujo 2: Selección de Receta
		var recipe string

		recipeForm := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Selecciona una receta de producción:").
					Options(
						huh.NewOption("💻 MERN Stack (React + Express + MongoDB)", "mern"),
						huh.NewOption("🚀 PERN Stack (React + Express + PostgreSQL)", "pern"),
						huh.NewOption("⚡ SaaS Starter (Next.js + Drizzle + Better-Auth + Stripe)", "saas"),
						huh.NewOption("🐍 FastAPI + React SPA", "fastapi_react"),
					).
					Value(&recipe),
			),
		)

		if err := recipeForm.Run(); err != nil {
			fmt.Println("Error ejecutando formulario:", err)
			os.Exit(1)
		}

		config := scaffold.ScaffoldConfig{
			ProjectName: projectName,
			Recipe:      recipe,
		}

		fmt.Printf("\n🚀 Inicializando '%s' con receta '%s'...\n", projectName, recipe)

		runInitSteps(projectName, config)

		fmt.Printf("\n✓ ¡Proyecto '%s' inicializado con éxito!\n", projectName)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVarP(&autoApprove, "default", "d", false, "Inicializar proyecto con configuración por defecto")
}
