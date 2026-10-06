package injector

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BlasVernazza06/koko-cli/internal/config"
)

func createDummyProject(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "koko-test-project-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create dummy koko.config.json
	kokoCfg := config.KokoConfig{
		Schema: "https://koko-cli.dev/schema.json",
		Project: config.ProjectInfo{
			Name:       "dummy-app",
			CLIVersion: "v1.0.0",
		},
		Architecture: config.ArchitectureInfo{
			Layout:         "monorepo",
			PackageManager: "pnpm",
		},
		Stack: config.StackInfo{
			Frontend: &config.FrontendInfo{Framework: "next", Language: "typescript"},
			Backend:  &config.BackendInfo{Framework: "express", Language: "typescript"},
			Database: &config.DatabaseInfo{Provider: "postgres", ORM: "drizzle"},
		},
	}

	data, err := json.MarshalIndent(kokoCfg, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal dummy config: %v", err)
	}

	if err := os.WriteFile(filepath.Join(tmpDir, "koko.config.json"), data, 0644); err != nil {
		t.Fatalf("Failed to write koko.config.json: %v", err)
	}

	// Create dummy apps/web and apps/api package.json
	_ = os.MkdirAll(filepath.Join(tmpDir, "apps", "web"), 0755)
	_ = os.MkdirAll(filepath.Join(tmpDir, "apps", "api"), 0755)

	_ = os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(`{"name": "dummy-root"}`), 0644)
	_ = os.WriteFile(filepath.Join(tmpDir, "apps/web/package.json"), []byte(`{"name": "@dummy-app/web", "dependencies": {}}`), 0644)
	_ = os.WriteFile(filepath.Join(tmpDir, "apps/api/package.json"), []byte(`{"name": "@dummy-app/api", "dependencies": {}}`), 0644)
	_ = os.WriteFile(filepath.Join(tmpDir, ".env"), []byte("# Environment Variables\n"), 0644)
	_ = os.WriteFile(filepath.Join(tmpDir, ".env.example"), []byte("# Example Variables\n"), 0644)

	return tmpDir
}

func TestAddAddon_Stripe(t *testing.T) {
	projectDir := createDummyProject(t)
	defer os.RemoveAll(projectDir)

	if err := AddAddon(projectDir, "stripe"); err != nil {
		t.Fatalf("AddAddon(stripe) failed: %v", err)
	}

	// 1. Verify koko.config.json updated
	data, err := os.ReadFile(filepath.Join(projectDir, "koko.config.json"))
	if err != nil {
		t.Fatalf("Failed to read koko.config.json: %v", err)
	}
	var cfg config.KokoConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("Failed to unmarshal updated config: %v", err)
	}
	if cfg.Features.Payments == nil || cfg.Features.Payments.Provider != "stripe" {
		t.Errorf("Expected Payments.Provider = stripe in config, got %+v", cfg.Features.Payments)
	}

	// 2. Verify apps/web/package.json contains stripe
	webPkgData, err := os.ReadFile(filepath.Join(projectDir, "apps/web/package.json"))
	if err != nil {
		t.Fatalf("Failed to read apps/web/package.json: %v", err)
	}
	if !strings.Contains(string(webPkgData), `"stripe"`) {
		t.Errorf("Expected 'stripe' in apps/web/package.json, got: %s", string(webPkgData))
	}

	// 3. Verify .env contains STRIPE_SECRET_KEY
	envData, err := os.ReadFile(filepath.Join(projectDir, ".env"))
	if err != nil {
		t.Fatalf("Failed to read .env: %v", err)
	}
	if !strings.Contains(string(envData), "STRIPE_SECRET_KEY") {
		t.Errorf("Expected STRIPE_SECRET_KEY in .env, got: %s", string(envData))
	}
}

func TestAddAddon_Idempotency(t *testing.T) {
	projectDir := createDummyProject(t)
	defer os.RemoveAll(projectDir)

	if err := AddAddon(projectDir, "resend"); err != nil {
		t.Fatalf("First AddAddon(resend) failed: %v", err)
	}

	// Second call should return an idempotency error
	err := AddAddon(projectDir, "resend")
	if err == nil {
		t.Fatalf("Expected error when adding resend twice, got nil")
	}
	if !strings.Contains(err.Error(), "already configured") {
		t.Errorf("Expected 'already configured' in error message, got: %v", err)
	}
}
