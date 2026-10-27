package cmd

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/BlasVernazza06/koko-cli/internal/features/skills"
	"github.com/BlasVernazza06/koko-cli/internal/tui/views"
)

func TestSkillsCmdRegistered(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"skills"})
	if err != nil {
		t.Fatalf("failed to find 'skills' command: %v", err)
	}
	if cmd == nil || cmd.Name() != "skills" {
		t.Fatalf("expected 'skills' command to be registered in rootCmd")
	}
	if !strings.Contains(cmd.Short, "Install and manage AI agent skills") {
		t.Errorf("unexpected Short description: %s", cmd.Short)
	}
}

func TestSkillsCmdDelegationToRunner(t *testing.T) {
	tmpDir := t.TempDir()

	origManager := skillsManager
	defer func() { skillsManager = origManager }()

	var (
		capturedDir  string
		capturedArgs []string
	)

	skillsManager = &skills.Manager{
		LookPath: func(file string) (string, error) {
			if file == "npx" {
				return "/fake/npx", nil
			}
			return "", errors.New("not found")
		},
		Exec: func(ctx context.Context, dir string, stdin io.Reader, stdout, stderr io.Writer, name string, args ...string) error {
			capturedDir = dir
			capturedArgs = args
			return nil
		},
	}

	resetFlags(rootCmd)
	cmd := rootCmd
	var outBuf, errBuf bytes.Buffer
	cmd.SetOut(&outBuf)
	cmd.SetErr(&errBuf)
	// Passing an explicit subcommand like "ui" triggers delegation
	cmd.SetArgs([]string{"skills", "--dir", tmpDir, "--yes", "--dry-run", "--agent", "cursor", "ui"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("skills command execution failed: %v", err)
	}

	if capturedDir != tmpDir {
		t.Errorf("capturedDir = %q, want %q", capturedDir, tmpDir)
	}

	expectedArgs := []string{"autoskills", "--yes", "--dry-run", "--agent", "cursor", "ui"}
	if !reflect.DeepEqual(capturedArgs, expectedArgs) {
		t.Errorf("capturedArgs = %v, want %v", capturedArgs, expectedArgs)
	}
}

func TestSkillsCmdNative_DryRun(t *testing.T) {
	tmpDir := t.TempDir()

	// Create dummy astro config
	if err := os.WriteFile(filepath.Join(tmpDir, "astro.config.mjs"), []byte("//"), 0644); err != nil {
		t.Fatalf("failed creating dummy file: %v", err)
	}

	resetFlags(rootCmd)
	cmd := rootCmd
	var outBuf, errBuf bytes.Buffer
	cmd.SetOut(&outBuf)
	cmd.SetErr(&errBuf)
	cmd.SetArgs([]string{"skills", "--dir", tmpDir, "--dry-run"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("skills --dry-run failed: %v", err)
	}

	// Make sure no .cursor/rules directory was created
	rulesDir := filepath.Join(tmpDir, ".cursor", "rules")
	if _, err := os.Stat(rulesDir); err == nil {
		t.Errorf("dry-run should not create .cursor/rules")
	}
}

func TestSkillsCmdNative_YesFlag(t *testing.T) {
	tmpDir := t.TempDir()

	// Write package.json with tailwindcss
	pkg := `{"dependencies": {"tailwindcss": "^3.0.0"}}`
	if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(pkg), 0644); err != nil {
		t.Fatalf("failed creating package.json: %v", err)
	}

	resetFlags(rootCmd)
	cmd := rootCmd
	var outBuf, errBuf bytes.Buffer
	cmd.SetOut(&outBuf)
	cmd.SetErr(&errBuf)
	cmd.SetArgs([]string{"skills", "--dir", tmpDir, "-y"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("skills -y failed: %v", err)
	}

	// Verify tailwindcss rule was injected
	ruleFile := filepath.Join(tmpDir, ".cursor", "rules", "tailwindcss.mdc")
	if _, err := os.Stat(ruleFile); os.IsNotExist(err) {
		t.Errorf("expected rule file at %s", ruleFile)
	}
}

func TestSkillsCmdNative_InteractiveTUI(t *testing.T) {
	tmpDir := t.TempDir()

	origOrch := skillsOrchestrator
	defer func() { skillsOrchestrator = origOrch }()

	orch, err := skills.NewOrchestrator()
	if err != nil {
		t.Fatalf("failed creating orchestrator: %v", err)
	}

	// Mock RunTUI returning vitest and nodejs
	orch.RunTUI = func(items []views.SkillItem, categories []string) ([]string, error) {
		return []string{"vitest", "nodejs"}, nil
	}
	skillsOrchestrator = orch

	resetFlags(rootCmd)
	cmd := rootCmd
	var outBuf, errBuf bytes.Buffer
	cmd.SetOut(&outBuf)
	cmd.SetErr(&errBuf)
	cmd.SetArgs([]string{"skills", "--dir", tmpDir})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("interactive execution failed: %v", err)
	}

	// Verify vitest.mdc and nodejs.mdc exist
	if _, err := os.Stat(filepath.Join(tmpDir, ".cursor", "rules", "vitest.mdc")); os.IsNotExist(err) {
		t.Errorf("expected vitest.mdc to be created")
	}
	if _, err := os.Stat(filepath.Join(tmpDir, ".cursor", "rules", "nodejs.mdc")); os.IsNotExist(err) {
		t.Errorf("expected nodejs.mdc to be created")
	}
}

func TestSkillsCmdNative_InvalidDir(t *testing.T) {
	resetFlags(rootCmd)
	cmd := rootCmd
	var outBuf, errBuf bytes.Buffer
	cmd.SetOut(&outBuf)
	cmd.SetErr(&errBuf)
	cmd.SetArgs([]string{"skills", "--dir", "c:\\non_existent_folder_987654321"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for non-existent directory, got nil")
	}
}
