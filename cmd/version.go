package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Muestra la versión de Koko CLI e información del sistema",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\n\033[90m┌\033[0m  \033[1mKoko CLI · Versión del sistema\033[0m")
		fmt.Println("\033[90m│\033[0m")
		fmt.Printf("\033[90m│\033[0m  \033[38;2;0;255;127m◆\033[0m  \033[1mVersión\033[0m     \033[90m·\033[0m  \033[38;2;167;139;250mv0.1.0\033[0m\n")
		fmt.Printf("\033[90m│\033[0m  \033[38;2;0;255;127m◆\033[0m  \033[1mOS/Arch\033[0m     \033[90m·\033[0m  \033[38;2;167;139;250m%s/%s\033[0m\n", runtime.GOOS, runtime.GOARCH)
		fmt.Println("\033[90m│\033[0m  \033[38;2;0;255;127m◆\033[0m  \033[1mBuild Date\033[0m  \033[90m·\033[0m  \033[38;2;167;139;250m2026-08-12\033[0m")
		fmt.Println("\033[90m│\033[0m")
		fmt.Println("\033[90m└\033[0m")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
