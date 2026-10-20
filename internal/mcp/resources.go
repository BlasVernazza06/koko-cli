package mcp

import (
	"context"
	"encoding/json"

	"github.com/BlasVernazza06/koko-cli/internal/catalog"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerResources(s *server.MCPServer) {
	// 1. Resource: koko://catalog/versions
	s.AddResource(
		mcp.NewResource("koko://catalog/versions",
			"Koko Catalog Versions",
			mcp.WithResourceDescription("The centralized dependency version catalog used across all Koko templates"),
			mcp.WithMIMEType("application/json"),
		),
		handleCatalogVersions,
	)
}

func handleCatalogVersions(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	bytes, err := json.MarshalIndent(catalog.DependencyVersions, "", "  ")
	if err != nil {
		return nil, err
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      "koko://catalog/versions",
			MIMEType: "application/json",
			Text:     string(bytes),
		},
	}, nil
}
