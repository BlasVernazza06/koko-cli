package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var autoApprove bool
var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize a new software project starter",
	Long:  `Run an interactive wizard to configure and generate a software project template.`,
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

		if autoApprove {
			fmt.Printf("🚀 Inicializando '%s' con configuración por defecto (SaaS Starter: Next.js + Go Fiber + PostgreSQL + Docker Compose)...\n", projectName)
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

		if setupMode == "quick" {
			var recipe string
			recipeForm := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Title("Selecciona una receta de producción:").
						Options(
							huh.NewOption("💻 Fullstack SaaS Starter (Next.js + Go Fiber + PostgreSQL + Docker Compose)", "saas"),
							huh.NewOption("⚡ API Moderna limpia (Fastify + Prisma + PostgreSQL)", "api"),
							huh.NewOption("🎨 Single Page App (React SPA + Vite + Tailwind CSS)", "spa"),
						).
						Value(&recipe),
				),
			)
			if err := recipeForm.Run(); err != nil {
				fmt.Println("Error ejecutando formulario:", err)
				os.Exit(1)
			}
			fmt.Printf("\n¡Inicializando '%s' con receta '%s'...\n", projectName, recipe)
			// Scaffolding behavior will be integrated here for v0.1.0
		} else {
			var frontend, backend, database string
			var docker, githubActions bool

			manualForm := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Title("Selecciona tu framework de Frontend:").
						Options(
							huh.NewOption("Next.js (App Router, TS)", "next"),
							huh.NewOption("React SPA (Vite, TS)", "react"),
							huh.NewOption("Ninguno", "none"),
						).
						Value(&frontend),
					huh.NewSelect[string]().
						Title("Selecciona tu framework de Backend:").
						Options(
							huh.NewOption("Go Fiber (REST API)", "fiber"),
							huh.NewOption("Node.js Express (TS)", "express"),
							huh.NewOption("Ninguno", "none"),
						).
						Value(&backend),
					huh.NewSelect[string]().
						Title("Selecciona tu ORM y Base de Datos:").
						Options(
							huh.NewOption("PostgreSQL + Prisma", "postgres_prisma"),
							huh.NewOption("PostgreSQL + SQLx (Go)", "postgres_sqlx"),
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

			fmt.Printf("\n¡Inicializando '%s' con configuración manual...\n", projectName)
			fmt.Printf("Frontend: %s, Backend: %s, Database: %s, Docker: %t, Github Actions: %t\n", frontend, backend, database, docker, githubActions)
			// Scaffolding behavior will be integrated here for v0.1.0
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVarP(&autoApprove, "default", "d", false, "Inicializar proyecto con configuración por defecto")
}
