package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	Version   = "v0.1.0"
	BuildDate = "2026-08-12"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version, OS/architecture, and build date of Koko-CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Koko-CLI %s\n", Version)
		fmt.Printf("OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
		fmt.Printf("Build Date: %s\n", BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
