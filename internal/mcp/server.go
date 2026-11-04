package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/server"
)

// ServerVersion define la versión del servidor MCP de Koko
const ServerVersion = "2.0.0"

// NewKokoServer inicializa un servidor MCP con todas las herramientas y recursos de Koko CLI.
func NewKokoServer() *server.MCPServer {
	s := server.NewMCPServer(
		"koko",
		ServerVersion,
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
	)

	// Registrar herramientas y recursos
	registerTools(s)
	registerResources(s)

	return s
}

// StartStdioServer inicia el servidor MCP escuchando a través del canal estándar stdio.
func StartStdioServer(ctx context.Context) error {
	s := NewKokoServer()
	if err := server.ServeStdio(s); err != nil {
		return fmt.Errorf("koko mcp stdio server error: %w", err)
	}
	return nil
}
