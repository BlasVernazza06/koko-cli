package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func resetFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		_ = f.Value.Set(f.DefValue)
		f.Changed = false
	})
	for _, sub := range cmd.Commands() {
		resetFlags(sub)
	}
}

func TestE2ECLIInitDefaultRecipe(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current wd: %v", err)
	}
	defer func() { _ = os.Chdir(originalWd) }()

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to chdir to temp dir: %v", err)
	}

	projectName := "e2e-saas-app"
	resetFlags(rootCmd)
	cmd := rootCmd
	var outBuf, errBuf bytes.Buffer
	cmd.SetOut(&outBuf)
	cmd.SetErr(&errBuf)
	cmd.SetArgs([]string{"init", projectName, "-d", "--git", "no"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("E2E CLI init failed: %v (stderr: %s)", err, errBuf.String())
	}

	targetDir := filepath.Join(tmpDir, projectName)
	expectedFiles := []string{
		filepath.Join(targetDir, "package.json"),
		filepath.Join(targetDir, "turbo.json"),
		filepath.Join(targetDir, "pnpm-workspace.yaml"),
		filepath.Join(targetDir, "koko.config.json"),
		filepath.Join(targetDir, "docker-compose.yml"),
		filepath.Join(targetDir, "apps", "web", "package.json"),
		filepath.Join(targetDir, "packages", "db", "package.json"),
		filepath.Join(targetDir, "packages", "auth", "package.json"),
	}

	for _, f := range expectedFiles {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			t.Errorf("E2E expected file %s does not exist on disk", f)
		}
	}
}

func TestE2ECLIInitManualFlags(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current wd: %v", err)
	}
	defer func() { _ = os.Chdir(originalWd) }()

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to chdir to temp dir: %v", err)
	}

	projectName := "e2e-manual-app"
	resetFlags(rootCmd)
	cmd := rootCmd
	var outBuf, errBuf bytes.Buffer
	cmd.SetOut(&outBuf)
	cmd.SetErr(&errBuf)
	cmd.SetArgs([]string{
		"init", projectName,
		"--frontend", "react",
		"--backend", "express",
		"--api", "trpc",
		"--database", "postgres",
		"--orm", "drizzle",
		"--auth", "better-auth",
		"--addons", "shadcn,stripe,docker",
		"--package-manager", "pnpm",
		"--git", "no",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("E2E CLI manual init failed: %v (stderr: %s)", err, errBuf.String())
	}

	targetDir := filepath.Join(tmpDir, projectName)
	expectedFiles := []string{
		filepath.Join(targetDir, "package.json"),
		filepath.Join(targetDir, "koko.config.json"),
		filepath.Join(targetDir, "docker-compose.yml"),
		filepath.Join(targetDir, "apps", "web", "package.json"),
		filepath.Join(targetDir, "apps", "api", "package.json"),
		filepath.Join(targetDir, "packages", "db", "package.json"),
		filepath.Join(targetDir, "packages", "ui", "package.json"),
	}

	for _, f := range expectedFiles {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			t.Errorf("E2E expected file %s does not exist on disk", f)
		}
	}
}

func TestE2ECLIDoctorCommand(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create a dummy project structure
	pkgJSON := `{"name": "test-doctor-app", "dependencies": {"next": "^14.0.0", "react": "^18.0.0"}}`
	if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(pkgJSON), 0644); err != nil {
		t.Fatalf("failed to write package.json: %v", err)
	}

	// 1. Run doctor with --dir flag (no fix)
	resetFlags(rootCmd)
	cmd := rootCmd
	var outBuf, errBuf bytes.Buffer
	cmd.SetOut(&outBuf)
	cmd.SetErr(&errBuf)
	cmd.SetArgs([]string{"doctor", "--dir", tmpDir})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("koko doctor --dir failed: %v (stderr: %s)", err, errBuf.String())
	}

	// koko.config.json should NOT exist yet
	configFile := filepath.Join(tmpDir, "koko.config.json")
	if _, err := os.Stat(configFile); !os.IsNotExist(err) {
		t.Errorf("koko.config.json should not exist before --fix")
	}

	// 2. Run doctor with -f / --fix and -d flag
	resetFlags(rootCmd)
	outBuf.Reset()
	errBuf.Reset()
	cmd.SetOut(&outBuf)
	cmd.SetErr(&errBuf)
	cmd.SetArgs([]string{"doctor", "-f", "-d", tmpDir})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("koko doctor -f -d failed: %v (stderr: %s)", err, errBuf.String())
	}

	// koko.config.json SHOULD exist now
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Errorf("koko.config.json was not generated after doctor -f")
	}

	// 3. Run doctor again: should show project in sync
	resetFlags(rootCmd)
	outBuf.Reset()
	errBuf.Reset()
	cmd.SetOut(&outBuf)
	cmd.SetErr(&errBuf)
	cmd.SetArgs([]string{"doctor", "--dir", tmpDir})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("koko doctor in-sync check failed: %v (stderr: %s)", err, errBuf.String())
	}
}
