package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var showVersion bool

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Claw-CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Claw-CLI v0.1.0")
	},
}

func init() {
	versionCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Print version number")
	rootCmd.AddCommand(versionCmd)
}
