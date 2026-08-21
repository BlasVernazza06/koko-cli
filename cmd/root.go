package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "koko",
	Short: "Koko CLI - Herramienta moderna de inicialización y scaffolding de proyectos",
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

