package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/BlasVernazza06/koko-cli/internal/mcp"
	"github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start the Koko Model Context Protocol (MCP) server over stdio",
	Long: `Starts the native Model Context Protocol (MCP) server for Koko CLI over stdio.
This enables AI coding agents (Claude Code, Cursor, Antigravity, Zed, Codex) to programmatically
scaffold projects, inject addons, validate compatibility matrices, and run health diagnostics.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		if err := mcp.StartStdioServer(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Error starting MCP server: %v\n", err)
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(mcpCmd)
}
