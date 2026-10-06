package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BlasVernazza06/koko-cli/internal/injector"
	"github.com/spf13/cobra"
)

var (
	addDirFlag string
)

var validAddons = map[string]string{
	"stripe":         "Payments (Stripe)",
	"polar":          "Payments (Polar.sh)",
	"resend":         "Email (Resend)",
	"brevo":          "Email (Brevo)",
	"shadcn":         "UI Components (shadcn/ui)",
	"lucide":         "Icons Pack (Lucide Icons)",
	"svgl":           "Brand Icons (SVGL)",
	"motion":         "Animations (Framer Motion)",
	"zod":            "Validation & Typing (Zod)",
	"docker":         "DevOps (Docker Compose)",
	"github_actions": "CI/CD (GitHub Actions)",
}

var addCmd = &cobra.Command{
	Use:   "add <addon>",
	Short: "Add a new service, dependency or tooling addon to an existing Koko project",
	Long: `Add a new service or addon (e.g. stripe, polar, resend, brevo, shadcn, lucide, zod, docker, github_actions)
to an existing Koko project by updating koko.config.json, installing dependencies and configuring environment variables.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addonArg := strings.ToLower(strings.TrimSpace(args[0]))

		description, isValid := validAddons[addonArg]
		if !isValid {
			fmt.Printf("\n\033[31m✗ Invalid addon: '%s'\033[0m\n\n", addonArg)
			fmt.Println("Available addons and services:")
			for key, desc := range validAddons {
				fmt.Printf("  • \033[1m%-16s\033[0m %s\n", key, desc)
			}
			fmt.Println()
			os.Exit(1)
		}

		targetDir := addDirFlag
		if targetDir == "" {
			targetDir = "."
		}

		configPath := filepath.Join(targetDir, "koko.config.json")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			fmt.Printf("\n\033[31m✗ Error: koko.config.json not found in '%s'\033[0m\n", targetDir)
			fmt.Println("  Please run 'koko add' inside an existing Koko project directory.")
			fmt.Println()
			os.Exit(1)
		}

		fmt.Printf("\n\033[90m┌\033[0m  \033[1mKoko Add · Injecting %s\033[0m\n", description)
		fmt.Println("\033[90m│\033[0m")

		if err := injector.AddAddon(targetDir, addonArg); err != nil {
			fmt.Printf("\033[90m│\033[0m  \033[31m✗ Injection error: %v\033[0m\n", err)
			fmt.Printf("\033[90m│\033[0m\n")
			fmt.Printf("\033[90m└\033[0m  \033[31m\033[1mFailed to add %s.\033[0m\n\n", addonArg)
			os.Exit(1)
		}

		fmt.Printf("\033[90m│\033[0m  \033[32m✓\033[0m  Updated koko.config.json\n")
		fmt.Printf("\033[90m│\033[0m  \033[32m✓\033[0m  Injected dependencies & environment variables\n")
		fmt.Printf("\033[90m│\033[0m\n")
		fmt.Printf("\033[90m└\033[0m  \033[32m\033[1mAddon '%s' successfully added to project!\033[0m\n\n", addonArg)
	},
}

func init() {
	addCmd.Flags().StringVarP(&addDirFlag, "dir", "d", ".", "Target project directory to modify")
	rootCmd.AddCommand(addCmd)
}