package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "claw",
	Short: "Claw-CLI is a modern project bootstrapper and scaffolder",
	Long: `A modern, high-performance command-line utility written in Go 
designed to eliminate friction in initializing, configuring, and structuring software projects.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Root flags and setup can be configured here
}
