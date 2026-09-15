package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "koko",
	Short:   "Koko CLI - Modern project initialization and scaffolding tool",
	Version: "v0.1.0",
	Run: func(cmd *cobra.Command, args []string) {
		RunTUI()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
