package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDBCmd_Overview(t *testing.T) {
	tmpDir := t.TempDir()

	configJSON := `{
		"architecture": {"layout": "monorepo", "packageManager": "pnpm"},
		"stack": {"database": {"provider": "postgres", "orm": "drizzle"}}
	}`
	_ = os.WriteFile(filepath.Join(tmpDir, "koko.config.json"), []byte(configJSON), 0644)

	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{"db", "--dir", tmpDir})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error executing koko db: %v", err)
	}

	result := out.String()
	if !strings.Contains(result, "Database Stack Overview") {
		t.Errorf("expected overview header in output: %s", result)
	}
	if !strings.Contains(result, "postgres") || !strings.Contains(result, "drizzle") {
		t.Errorf("expected stack details in output: %s", result)
	}
}

func TestDBCmd_SubcommandsDryRun(t *testing.T) {
	tmpDir := t.TempDir()

	configJSON := `{
		"architecture": {"layout": "standalone", "packageManager": "npm"},
		"stack": {"database": {"provider": "postgres", "orm": "prisma"}}
	}`
	_ = os.WriteFile(filepath.Join(tmpDir, "koko.config.json"), []byte(configJSON), 0644)

	subcommands := []string{"migrate", "push", "generate", "studio", "seed", "reset"}

	for _, sub := range subcommands {
		t.Run(sub, func(t *testing.T) {
			var out bytes.Buffer
			rootCmd.SetOut(&out)
			rootCmd.SetErr(&out)
			rootCmd.SetArgs([]string{"db", sub, "--dir", tmpDir, "--dry-run"})

			err := rootCmd.Execute()
			if err != nil {
				t.Fatalf("unexpected error executing koko db %s --dry-run: %v", sub, err)
			}

			result := out.String()
			if !strings.Contains(result, "Dry Run Analysis") {
				t.Errorf("expected dry run analysis header in output: %s", result)
			}
			if !strings.Contains(result, "npx prisma") {
				t.Errorf("expected npx prisma in dry-run output: %s", result)
			}
		})
	}
}

func TestDBCmd_PathCleanAndInvalidDir(t *testing.T) {
	// Probar que un directorio no existente devuelva error controlado sin pánico
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{"db", "--dir", "c:\\non_existent_dir_999999"})

	err := rootCmd.Execute()
	if err == nil {
		t.Errorf("expected error for non-existent directory, got nil")
	}
}

