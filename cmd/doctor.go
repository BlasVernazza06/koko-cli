package cmd

import (
	"github.com/spf13/cobra"
)

var doctorCommand = &cobra.Command{
	Use:   "doctor",
	Short: "Verify your project, looking for dependencies that are not declared",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
