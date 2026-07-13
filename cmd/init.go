package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/BlasVernazza06/koko-cli/internal/scaffold"
)

var autoApprove bool
var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize a new software project starter",
	Long:  `Run an interactive wizard to configure and generate a software project template.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := "my-project"
		if len(args) > 0 {
			projectName = args[0]
		} else if !autoApprove {
			nameForm := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("¿Cuál es el nombre de tu proyecto?").
						Placeholder("my-project").
						Value(&projectName),
				),
			)
			if err := nameForm.Run(); err != nil {
				fmt.Println("Error ejecutando formulario:", err)
				os.Exit(1)
			}
		}

		fmt.Println(`
			 __ _    __ _ __      __
			/  _ | |/ _' |\ \ /\ / /
			| (__| | (_| | \ V  V / 
			\___ |_|\__,_|  \_/\_/  v0.1.0
		`)

		// Flujo 1: Inicialización automática con configuración por defecto (Next + Express + Postgres + Docker)
		if autoApprove {
			fmt.Printf("🚀 Inicializando '%s' con configuración por defecto (Next.js, Express, pnpm, PostgreSQL, Drizzle, Turborepo, Docker)...\n", projectName)
			config := scaffold.ScaffoldConfig{
				ProjectName:   projectName,
				Frontend:      scaffold.FrontendNext,
				Backend:       scaffold.BackendExpress,
				Database:      scaffold.DatabasePostgres,
				Docker:        true,
				GithubActions: true,
			}
			err := scaffold.RunScaffold(projectName, config)
			if err != nil {
				fmt.Printf("❌ Error al inicializar el proyecto: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("✓ ¡Proyecto '%s' inicializado con éxito!\n", projectName)
			return
		}

		var setupMode string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("¿Cómo deseas inicializar tu proyecto?").
					Options(
						huh.NewOption("🚀 Setup Rápido (Recetas de producción)", "quick"),
						huh.NewOption("⚙️  Configuración Manual (Elegir stack paso a paso)", "manual"),
					).
					Value(&setupMode),
			),
		)

		err := form.Run()
		if err != nil {
			fmt.Println("Error ejecutando formulario:", err)
			os.Exit(1)
		}

		var config scaffold.ScaffoldConfig
		config.ProjectName = projectName

		// Flujo 2: Setup Rápido (Recetas)
		if setupMode == "quick" {
			var recipe string
			recipeForm := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Title("Selecciona una receta de producción:").
						Options(
							huh.NewOption("💻 Fullstack SaaS Starter (Next.js + Express + PostgreSQL + Docker)", "saas"),
							huh.NewOption("⚡ API Moderna limpia (Express + PostgreSQL + Docker)", "api"),
							huh.NewOption("🎨 Single Page App (React SPA + Vite + Tailwind CSS)", "spa"),
						).
						Value(&recipe),
				),
			)
			if err := recipeForm.Run(); err != nil {
				fmt.Println("Error ejecutando formulario:", err)
				os.Exit(1)
			}

			switch recipe {
			case "saas":
				config.Frontend = scaffold.FrontendNext
				config.Backend = scaffold.BackendExpress
				config.Database = scaffold.DatabasePostgres
				config.Docker = true
				config.GithubActions = true
			case "api":
				config.Frontend = scaffold.FrontendNone
				config.Backend = scaffold.BackendExpress
				config.Database = scaffold.DatabasePostgres
				config.Docker = true
				config.GithubActions = false
			case "spa":
				config.Frontend = scaffold.FrontendReact
				config.Backend = scaffold.BackendNone
				config.Database = scaffold.DatabaseNone
				config.Docker = false
				config.GithubActions = false
			}

			fmt.Printf("\n¡Inicializando '%s' con receta '%s'...\n", projectName, recipe)

		} else {
			// Flujo 3: Configuración Manual paso a paso
			var frontend, backend, database string
			var docker, githubActions bool

			manualForm := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Title("Selecciona tu framework de Frontend:").
						Options(
							huh.NewOption("Next.js (App Router, TS)", "next"),
							huh.NewOption("React SPA (Vite, TS)", "react"),
							huh.NewOption("Vue.js (Vite, TS)", "vue"),
							huh.NewOption("Ninguno", "none"),
						).
						Value(&frontend),
					huh.NewSelect[string]().
						Title("Selecciona tu framework de Backend:").
						Options(
							huh.NewOption("Go Fiber (REST API)", "fiber"),
							huh.NewOption("Node.js Express (TS)", "express"),
							huh.NewOption("Hono (Node.js Server)", "hono"),
							huh.NewOption("Ninguno", "none"),
						).
						Value(&backend),
					huh.NewSelect[string]().
						Title("Selecciona tu Base de Datos:").
						Options(
							huh.NewOption("PostgreSQL", "postgres"),
							huh.NewOption("MySQL", "mysql"),
							huh.NewOption("MongoDB", "mongodb"),
							huh.NewOption("Ninguno", "none"),
						).
						Value(&database),
				),
				huh.NewGroup(
					huh.NewConfirm().
						Title("¿Configurar entorno de desarrollo local con Docker Compose?").
						Value(&docker),
					huh.NewConfirm().
						Title("¿Configurar Github Actions para CI/CD?").
						Value(&githubActions),
				),
			)

			if err := manualForm.Run(); err != nil {
				fmt.Println("Error ejecutando formulario:", err)
				os.Exit(1)
			}

			// Mapeamos los strings a los tipos fuertemente tipados de scaffold
			config.Frontend = scaffold.FrontendType(frontend)
			config.Backend = scaffold.BackendType(backend)
			config.Database = scaffold.DatabaseType(database)
			config.Docker = docker
			config.GithubActions = githubActions

			fmt.Printf("\n¡Inicializando '%s' con configuración manual...\n", projectName)
		}

		// Ejecutamos el scaffolding físico en el disco
		err = scaffold.RunScaffold(projectName, config)
		if err != nil {
			fmt.Printf("❌ Error al generar estructura de archivos: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ ¡Proyecto '%s' inicializado con éxito!\n", projectName)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVarP(&autoApprove, "default", "d", false, "Inicializar proyecto con configuración por defecto")
}
