package mcp

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewKokoServer(t *testing.T) {
	s := NewKokoServer()
	require.NotNil(t, s)
}

func TestHandleGetOptions(t *testing.T) {
	ctx := context.Background()
	req := mcp.CallToolRequest{}

	res, err := handleGetOptions(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.False(t, res.IsError)

	var data map[string]interface{}
	err = json.Unmarshal([]byte(res.Content[0].(mcp.TextContent).Text), &data)
	require.NoError(t, err)
	assert.Contains(t, data, "recipes")
	assert.Contains(t, data, "frontends")
	assert.Contains(t, data, "addons")
}

func TestHandleValidateStack(t *testing.T) {
	ctx := context.Background()

	// Valido: Next.js + Hono + Postgres + Drizzle + Better-Auth
	validReq := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"frontend": "next",
				"backend":  "hono",
				"database": "postgres",
				"orm":      "drizzle",
				"auth":     "better-auth",
			},
		},
	}
	res, err := handleValidateStack(ctx, validReq)
	require.NoError(t, err)
	assert.False(t, res.IsError)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "compatible")
}

func TestHandleCatalogVersionsResource(t *testing.T) {
	ctx := context.Background()
	req := mcp.ReadResourceRequest{
		Params: mcp.ReadResourceParams{
			URI: "koko://catalog/versions",
		},
	}

	contents, err := handleCatalogVersions(ctx, req)
	require.NoError(t, err)
	require.Len(t, contents, 1)

	textContent, ok := contents[0].(mcp.TextResourceContents)
	require.True(t, ok)
	assert.Equal(t, "koko://catalog/versions", textContent.URI)
	assert.Contains(t, textContent.Text, "next")
	assert.Contains(t, textContent.Text, "react")
}

func TestHandleDoctorTool(t *testing.T) {
	ctx := context.Background()
	tmpDir, err := os.MkdirTemp("", "koko-mcp-doctor-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"projectDirectory": tmpDir,
				"autoFix":          false,
			},
		},
	}

	res, err := handleDoctor(ctx, req)
	require.NoError(t, err)
	assert.False(t, res.IsError)

	var docSummary map[string]interface{}
	err = json.Unmarshal([]byte(res.Content[0].(mcp.TextContent).Text), &docSummary)
	require.NoError(t, err)
	assert.Equal(t, false, docSummary["configExists"])
}

func TestHandleSetAndGetDefaultDirectory(t *testing.T) {
	ctx := context.Background()
	tmpDir, err := os.MkdirTemp("", "koko-mcp-target-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Set
	setReq := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"directory": tmpDir,
			},
		},
	}
	setRes, err := handleSetDefaultDirectory(ctx, setReq)
	require.NoError(t, err)
	assert.False(t, setRes.IsError)
	assert.Contains(t, setRes.Content[0].(mcp.TextContent).Text, tmpDir)

	// Get
	getReq := mcp.CallToolRequest{}
	getRes, err := handleGetDefaultDirectory(ctx, getReq)
	require.NoError(t, err)
	assert.False(t, getRes.IsError)

	var resData map[string]interface{}
	err = json.Unmarshal([]byte(getRes.Content[0].(mcp.TextContent).Text), &resData)
	require.NoError(t, err)
	assert.Equal(t, tmpDir, resData["configuredDirectory"])
}

func TestHandleInitProjectValidation(t *testing.T) {
	ctx := context.Background()
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}

	res, err := handleInitProject(ctx, req)
	require.NoError(t, err)
	assert.True(t, res.IsError)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "projectName is required")
}

func TestHandleInitProjectSuccess(t *testing.T) {
	ctx := context.Background()
	tmpDir, err := os.MkdirTemp("", "koko-mcp-init-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"projectName":     "test-mcp-app",
				"recipe":          "saas-next",
				"targetDirectory": tmpDir,
				"initGit":         false,
			},
		},
	}

	res, err := handleInitProject(ctx, req)
	require.NoError(t, err)
	assert.False(t, res.IsError)

	createdPath := filepath.Join(tmpDir, "test-mcp-app")
	assert.DirExists(t, createdPath)
	assert.FileExists(t, filepath.Join(createdPath, "koko.config.json"))
}
