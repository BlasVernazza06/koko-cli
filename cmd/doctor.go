package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BlasVernazza06/koko-cli/internal/features/doctor"
	"github.com/spf13/cobra"
)

var (
	fixFlag bool
	dirFlag string
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose your project stack, environment tools, ports, and configuration",
	Long: `koko doctor is an all-in-one local health check and observability diagnostics tool.
It inspects your local environment tools (Node.js, Package Manager, Git, Docker),
stack architecture, environment variables (.env vs .env.example), database status,
port availability, AI skills context, and manifest synchronization.`,
	Run: func(cmd *cobra.Command, args []string) {
		rawDir := dirFlag
		if rawDir == "" {
			rawDir = "."
		}
		targetDir := filepath.Clean(rawDir)
		isFix := fixFlag

		fmt.Println("\n\033[90m┌\033[0m  \033[1mKoko Doctor · Comprehensive Project Health Check\033[0m")
		fmt.Printf("\033[90m│\033[0m  Directory: \033[38;2;167;139;250m%s\033[0m\n", targetDir)
		fmt.Println("\033[90m│\033[0m")

		result, err := doctor.RunDiagnostics(targetDir)
		if err != nil {
			fmt.Printf("\033[90m│\033[0m  \033[31m✗ Diagnostics error: %v\033[0m\n", err)
			fmt.Println("\033[90m│\033[0m")
			fmt.Printf("\033[90m└\033[0m  \033[31m\033[1mFailed to analyze project.\033[0m\n\n")
			os.Exit(1)
		}

		// Helper para renderizar items de salud
		printHealthItem := func(item doctor.HealthItem) {
			symbol := "\033[32m✔\033[0m"
			switch item.Status {
			case doctor.StatusWarn:
				symbol = "\033[33m⚠\033[0m"
			case doctor.StatusFail:
				symbol = "\033[31m✖\033[0m"
			case doctor.StatusInfo:
				symbol = "\033[38;2;167;139;250m◆\033[0m"
			}
			msg := item.Message
			if item.Detail != "" {
				msg = fmt.Sprintf("%s \033[90m(%s)\033[0m", msg, item.Detail)
			}
			fmt.Printf("\033[90m│\033[0m    %s  \033[1m%-28s\033[0m \033[90m·\033[0m  %s\n", symbol, item.Name, msg)
		}

		// =========================================================================
		// [1/5] Entorno & Herramientas
		// =========================================================================
		fmt.Println("\033[90m│\033[0m  \033[1m[1/5] Entorno & Herramientas Locales\033[0m")
		if result.HealthReport != nil {
			for _, item := range result.HealthReport.Environment {
				printHealthItem(item)
			}
		}

		// =========================================================================
		// [2/5] Stack & Arquitectura
		// =========================================================================
		fmt.Println("\033[90m│\033[0m")
		fmt.Println("\033[90m│\033[0m  \033[1m[2/5] Stack & Arquitectura Detectada\033[0m")
		fmt.Printf("\033[90m│\033[0m    \033[38;2;167;139;250m◆\033[0m  \033[1m%-28s\033[0m \033[90m·\033[0m  %s\n", "Layout", result.DetectedConfig.Architecture.Layout)
		fmt.Printf("\033[90m│\033[0m    \033[38;2;167;139;250m◆\033[0m  \033[1m%-28s\033[0m \033[90m·\033[0m  %s\n", "Package Manager", result.DetectedConfig.Architecture.PackageManager)

		if result.DetectedConfig.Stack.Frontend != nil {
			feInfo := fmt.Sprintf("%s (%s)", result.DetectedConfig.Stack.Frontend.Framework, result.DetectedConfig.Stack.Frontend.Language)
			if result.DetectedConfig.Stack.Frontend.Styling != "" && result.DetectedConfig.Stack.Frontend.Styling != "none" {
				feInfo += fmt.Sprintf(" + %s", result.DetectedConfig.Stack.Frontend.Styling)
			}
			fmt.Printf("\033[90m│\033[0m    \033[38;2;167;139;250m◆\033[0m  \033[1m%-28s\033[0m \033[90m·\033[0m  %s\n", "Frontend", feInfo)
		}
		if result.DetectedConfig.Stack.Backend != nil {
			beInfo := fmt.Sprintf("%s (%s)", result.DetectedConfig.Stack.Backend.Framework, result.DetectedConfig.Stack.Backend.Language)
			fmt.Printf("\033[90m│\033[0m    \033[38;2;167;139;250m◆\033[0m  \033[1m%-28s\033[0m \033[90m·\033[0m  %s\n", "Backend", beInfo)
		}
		if result.DetectedConfig.Stack.Database != nil {
			dbInfo := fmt.Sprintf("%s (ORM: %s)", result.DetectedConfig.Stack.Database.Provider, result.DetectedConfig.Stack.Database.ORM)
			fmt.Printf("\033[90m│\033[0m    \033[38;2;167;139;250m◆\033[0m  \033[1m%-28s\033[0m \033[90m·\033[0m  %s\n", "Database", dbInfo)
		}

		// =========================================================================
		// [3/5] Configuración & Variables (.env)
		// =========================================================================
		fmt.Println("\033[90m│\033[0m")
		fmt.Println("\033[90m│\033[0m  \033[1m[3/5] Configuración & Variables (.env)\033[0m")
		if result.HealthReport != nil {
			for _, item := range result.HealthReport.EnvVars {
				printHealthItem(item)
			}
		}

		// =========================================================================
		// [4/5] Base de Datos & Puertos Locales
		// =========================================================================
		fmt.Println("\033[90m│\033[0m")
		fmt.Println("\033[90m│\033[0m  \033[1m[4/5] Base de Datos & Puertos Locales\033[0m")
		if result.HealthReport != nil {
			for _, item := range result.HealthReport.DatabaseAndPorts {
				printHealthItem(item)
			}
			for _, item := range result.HealthReport.Skills {
				printHealthItem(item)
			}
		}

		// =========================================================================
		// [5/5] Sincronización de Manifiesto (koko.config.json)
		// =========================================================================
		fmt.Println("\033[90m│\033[0m")
		fmt.Println("\033[90m│\033[0m  \033[1m[5/5] Sincronización de Manifiesto (koko.config.json)\033[0m")
		if !result.ConfigExists {
			fmt.Printf("\033[90m│\033[0m    \033[33m⚠\033[0m  \033[1m%-28s\033[0m \033[90m·\033[0m  \033[38;2;255;184;108mNot found (will be generated with --fix)\033[0m\n", "Manifest Status")
		} else if len(result.Differences) == 0 {
			fmt.Printf("\033[90m│\033[0m    \033[32m✔\033[0m  \033[1m%-28s\033[0m \033[90m·\033[0m  \033[32mFully in sync with detected stack\033[0m\n", "Manifest Status")
		} else {
			fmt.Printf("\033[90m│\033[0m    \033[33m⚠\033[0m  \033[1m%-28s\033[0m \033[90m·\033[0m  \033[33m%d discrepancy(ies) detected\033[0m\n", "Manifest Status", len(result.Differences))
			for _, diff := range result.Differences {
				switch diff.Type {
				case doctor.DiffAdded:
					fmt.Printf("\033[90m│\033[0m        \033[32m+ [ADDED]\033[0m    %s > %s: \033[32m%s\033[0m\n", diff.Section, diff.Field, diff.NewValue)
				case doctor.DiffModified:
					fmt.Printf("\033[90m│\033[0m        \033[33m~ [MODIFIED]\033[0m %s > %s: \033[90m%s\033[0m → \033[38;2;167;139;250m%s\033[0m\n", diff.Section, diff.Field, diff.OldValue, diff.NewValue)
				case doctor.DiffRemoved:
					fmt.Printf("\033[90m│\033[0m        \033[31m- [REMOVED]\033[0m  %s > %s: \033[31m%s\033[0m\n", diff.Section, diff.Field, diff.OldValue)
				}
			}
		}

		// Ejecución de --fix si fue solicitada
		if isFix {
			fmt.Println("\033[90m│\033[0m")
			if err := doctor.ApplyFixes(targetDir, result.DetectedConfig, result.ExistingConfig); err != nil {
				fmt.Printf("\033[90m└\033[0m  \033[31m✗ Failed to apply fixes: %v\033[0m\n\n", err)
				os.Exit(1)
			}
			fmt.Printf("\033[90m└\033[0m  \033[32m\033[1m✔ koko.config.json and missing .env files were successfully updated!\033[0m\n\n")
			return
		}

		// Resumen y Acciones Recomendadas
		fmt.Println("\033[90m│\033[0m")
		var allActions []string
		if result.HealthReport != nil {
			allActions = append(allActions, result.HealthReport.RecommendedActions...)
		}
		if !result.ConfigExists || len(result.Differences) > 0 {
			allActions = append(allActions, "koko doctor --fix")
		}

		if len(allActions) > 0 {
			fmt.Println("\033[90m│\033[0m  \033[1m💡 Acciones Recomendadas para Resolver Alertas:\033[0m")
			seen := make(map[string]bool)
			for _, act := range allActions {
				act = strings.TrimSpace(act)
				if act != "" && !seen[act] {
					seen[act] = true
					fmt.Printf("\033[90m│\033[0m    \033[32m•\033[0m \033[38;2;167;139;250m%s\033[0m\n", act)
				}
			}
			fmt.Println("\033[90m│\033[0m")
			fmt.Println("\033[90m└\033[0m  \033[33mRun \"koko doctor --fix\" (or -f) to automatically sync manifest and create .env files.\033[0m")
			fmt.Println()
		} else {
			fmt.Println("\033[90m└\033[0m  \033[32m\033[1m✔ All systems nominal! Your local environment and project are ready.\033[0m")
			fmt.Println()
		}
	},
}

func init() {
	doctorCmd.Flags().BoolVarP(&fixFlag, "fix", "f", false, "Apply detected changes to koko.config.json and create missing .env files")
	doctorCmd.Flags().StringVarP(&dirFlag, "dir", "d", ".", "Target project directory to analyze")
	rootCmd.AddCommand(doctorCmd)
}
