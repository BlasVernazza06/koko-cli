package cmd

import (
	"fmt"
	"os"

	"github.com/BlasVernazza06/koko-cli/internal/doctor"
	"github.com/spf13/cobra"
)

var (
	fixFlag bool
	dirFlag string
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose your project stack, detect discrepancies and update koko.config.json",
	Run: func(cmd *cobra.Command, args []string) {
		targetDir, _ := cmd.Flags().GetString("dir")
		if targetDir == "" {
			targetDir = "."
		}
		isFix, _ := cmd.Flags().GetBool("fix")

		fmt.Println("\n\033[90m┌\033[0m  \033[1mKoko Doctor · Project Diagnostics\033[0m")
		fmt.Println("\033[90m│\033[0m")

		result, err := doctor.RunDiagnostics(targetDir)
		if err != nil {
			fmt.Printf("\033[90m│\033[0m  \033[31m✗ Diagnostics error: %v\033[0m\n", err)
			fmt.Printf("\033[90m│\033[0m\n")
			fmt.Printf("\033[90m└\033[0m  \033[31m\033[1mFailed to analyze project.\033[0m\n\n")
			os.Exit(1)
		}

		if !result.ConfigExists {
			fmt.Printf("\033[90m│\033[0m  \033[38;2;255;184;108m!\033[0m  \033[1mManifest\033[0m          \033[90m·\033[0m  \033[38;2;255;184;108mkoko.config.json not found (will be generated)\033[0m\n")
		} else {
			fmt.Printf("\033[90m│\033[0m  \033[38;2;0;255;127m◆\033[0m  \033[1mManifest\033[0m          \033[90m·\033[0m  \033[38;2;167;139;250mkoko.config.json detected\033[0m\n")
		}

		fmt.Printf("\033[90m│\033[0m  \033[38;2;0;255;127m◆\033[0m  \033[1mLayout\033[0m            \033[90m·\033[0m  \033[38;2;167;139;250m%s\033[0m\n", result.DetectedConfig.Architecture.Layout)
		fmt.Printf("\033[90m│\033[0m  \033[38;2;0;255;127m◆\033[0m  \033[1mPackage Manager\033[0m   \033[90m·\033[0m  \033[38;2;167;139;250m%s\033[0m\n", result.DetectedConfig.Architecture.PackageManager)

		if result.DetectedConfig.Stack.Frontend != nil {
			fmt.Printf("\033[90m│\033[0m  \033[38;2;0;255;127m◆\033[0m  \033[1mFrontend\033[0m          \033[90m·\033[0m  \033[38;2;167;139;250m%s (%s)\033[0m\n",
				result.DetectedConfig.Stack.Frontend.Framework,
				result.DetectedConfig.Stack.Frontend.Language)
		}
		if result.DetectedConfig.Stack.Backend != nil {
			fmt.Printf("\033[90m│\033[0m  \033[38;2;0;255;127m◆\033[0m  \033[1mBackend\033[0m           \033[90m·\033[0m  \033[38;2;167;139;250m%s (%s)\033[0m\n",
				result.DetectedConfig.Stack.Backend.Framework,
				result.DetectedConfig.Stack.Backend.Language)
		}
		if result.DetectedConfig.Stack.Database != nil {
			fmt.Printf("\033[90m│\033[0m  \033[38;2;0;255;127m◆\033[0m  \033[1mDatabase\033[0m          \033[90m·\033[0m  \033[38;2;167;139;250m%s (ORM: %s)\033[0m\n",
				result.DetectedConfig.Stack.Database.Provider,
				result.DetectedConfig.Stack.Database.ORM)
		}

		if len(result.Differences) == 0 {
			fmt.Printf("\033[90m│\033[0m\n")
			fmt.Printf("\033[90m└\033[0m  \033[32m\033[1m✓ Project is fully in sync with koko.config.json! No changes needed.\033[0m\n\n")
			return
		}

		fmt.Printf("\033[90m│\033[0m\n")
		fmt.Printf("\033[90m│\033[0m  \033[1mDetected Differences (%d):\033[0m\n", len(result.Differences))
		for _, diff := range result.Differences {
			switch diff.Type {
			case doctor.DiffAdded:
				fmt.Printf("\033[90m│\033[0m    \033[32m+ [ADDED]\033[0m    %s > %s: \033[32m%s\033[0m\n", diff.Section, diff.Field, diff.NewValue)
			case doctor.DiffModified:
				fmt.Printf("\033[90m│\033[0m    \033[33m~ [MODIFIED]\033[0m %s > %s: \033[90m%s\033[0m → \033[38;2;167;139;250m%s\033[0m\n", diff.Section, diff.Field, diff.OldValue, diff.NewValue)
			case doctor.DiffRemoved:
				fmt.Printf("\033[90m│\033[0m    \033[31m- [REMOVED]\033[0m  %s > %s: \033[31m%s\033[0m\n", diff.Section, diff.Field, diff.OldValue)
			}
		}

		fmt.Printf("\033[90m│\033[0m\n")
		if isFix {
			if err := doctor.ApplyFixes(targetDir, result.DetectedConfig, result.ExistingConfig); err != nil {
				fmt.Printf("\033[90m└\033[0m  \033[31m✗ Failed to update koko.config.json: %v\033[0m\n\n", err)
				os.Exit(1)
			}
			fmt.Printf("\033[90m└\033[0m  \033[32m\033[1m✓ koko.config.json successfully updated with detected changes!\033[0m\n\n")
		} else {
			fmt.Printf("\033[90m└\033[0m  \033[33mRun \"koko doctor --fix\" (or -f) to update koko.config.json with these changes.\033[0m\n\n")
		}
	},
}

func init() {
	doctorCmd.Flags().BoolVarP(&fixFlag, "fix", "f", false, "Apply detected changes to koko.config.json automatically")
	doctorCmd.Flags().StringVarP(&dirFlag, "dir", "d", ".", "Target project directory to analyze")
	rootCmd.AddCommand(doctorCmd)
}
