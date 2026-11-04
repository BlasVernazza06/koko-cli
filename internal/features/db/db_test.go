package db

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDetectProjectDBContext_FromKokoConfig(t *testing.T) {
	tmpDir := t.TempDir()

	configJSON := `{
		"schema": "https://koko-cli.dev/schema.json",
		"architecture": {
			"layout": "monorepo",
			"packageManager": "pnpm"
		},
		"stack": {
			"database": {
				"provider": "postgres",
				"orm": "drizzle"
			}
		}
	}`

	if err := os.WriteFile(filepath.Join(tmpDir, "koko.config.json"), []byte(configJSON), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	pCtx, err := DetectProjectDBContext(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !pCtx.ConfigExists {
		t.Errorf("expected ConfigExists = true")
	}
	if pCtx.Layout != "monorepo" {
		t.Errorf("expected layout monorepo, got %s", pCtx.Layout)
	}
	if pCtx.PackageManager != "pnpm" {
		t.Errorf("expected pm pnpm, got %s", pCtx.PackageManager)
	}
	if pCtx.Provider != "postgres" {
		t.Errorf("expected provider postgres, got %s", pCtx.Provider)
	}
	if pCtx.ORM != "drizzle" {
		t.Errorf("expected orm drizzle, got %s", pCtx.ORM)
	}
}

func TestResolveCommand_PackageJSONScriptsPriority(t *testing.T) {
	pCtx := &ProjectDBContext{
		ProjectDir:     "/test/project",
		Layout:         "monorepo",
		PackageManager: "pnpm",
		Provider:       "postgres",
		ORM:            "drizzle",
		RootScripts: map[string]string{
			"db:migrate":  "pnpm --filter @repo/db db:migrate",
			"db:push":     "pnpm --filter @repo/db db:push",
			"db:studio":   "pnpm --filter @repo/db db:studio",
			"db:generate": "pnpm --filter @repo/db db:generate",
			"db:seed":     "pnpm --filter @repo/db db:seed",
			"db:reset":    "pnpm --filter @repo/db db:reset",
		},
	}

	tests := []struct {
		action    Action
		extraArgs []string
		wantExe   string
		wantArgs  []string
	}{
		{
			action:   ActionMigrate,
			wantExe:  "pnpm",
			wantArgs: []string{"run", "db:migrate"},
		},
		{
			action:    ActionMigrate,
			extraArgs: []string{"--name", "add_users"},
			wantExe:   "pnpm",
			wantArgs:  []string{"run", "db:migrate", "--", "--name", "add_users"},
		},
		{
			action:   ActionPush,
			wantExe:  "pnpm",
			wantArgs: []string{"run", "db:push"},
		},
		{
			action:   ActionStudio,
			wantExe:  "pnpm",
			wantArgs: []string{"run", "db:studio"},
		},
		{
			action:   ActionSeed,
			wantExe:  "pnpm",
			wantArgs: []string{"run", "db:seed"},
		},
		{
			action:   ActionReset,
			wantExe:  "pnpm",
			wantArgs: []string{"run", "db:reset"},
		},
	}

	for _, tt := range tests {
		t.Run(string(tt.action), func(t *testing.T) {
			cmd, err := ResolveCommand(pCtx, tt.action, tt.extraArgs)
			if err != nil {
				t.Fatalf("unexpected error resolving %s: %v", tt.action, err)
			}
			if cmd.Executable != tt.wantExe {
				t.Errorf("exe = %s, want %s", cmd.Executable, tt.wantExe)
			}
			if len(cmd.Args) != len(tt.wantArgs) {
				t.Fatalf("args len = %d (%v), want %d (%v)", len(cmd.Args), cmd.Args, len(tt.wantArgs), tt.wantArgs)
			}
			for i, a := range cmd.Args {
				if a != tt.wantArgs[i] {
					t.Errorf("arg[%d] = %s, want %s", i, a, tt.wantArgs[i])
				}
			}
		})
	}
}

func TestResolveCommand_NativeDrizzleFallback(t *testing.T) {
	pCtx := &ProjectDBContext{
		ProjectDir:     "/test/standalone",
		Layout:         "standalone",
		PackageManager: "pnpm",
		Provider:       "postgres",
		ORM:            "drizzle",
		RootScripts:    map[string]string{}, // Sin scripts
	}

	cmd, err := ResolveCommand(pCtx, ActionPush, []string{"--force"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cmd.Executable != "pnpm" {
		t.Errorf("expected pnpm, got %s", cmd.Executable)
	}
	expectedArgs := []string{"exec", "drizzle-kit", "push", "--force"}
	if strings.Join(cmd.Args, " ") != strings.Join(expectedArgs, " ") {
		t.Errorf("args = %v, want %v", cmd.Args, expectedArgs)
	}

	// Acción no soportada nativamente por Drizzle sin script
	_, err = ResolveCommand(pCtx, ActionSeed, nil)
	if err == nil {
		t.Errorf("expected error for seed on drizzle without script, got nil")
	}
}

func TestResolveCommand_NativePrismaFallback(t *testing.T) {
	pCtx := &ProjectDBContext{
		ProjectDir:     "/test/standalone",
		Layout:         "standalone",
		PackageManager: "npm",
		Provider:       "postgres",
		ORM:            "prisma",
		RootScripts:    map[string]string{},
	}

	tests := []struct {
		action   Action
		wantArgs []string
	}{
		{action: ActionMigrate, wantArgs: []string{"prisma", "migrate", "dev"}},
		{action: ActionPush, wantArgs: []string{"prisma", "db", "push"}},
		{action: ActionGenerate, wantArgs: []string{"prisma", "generate"}},
		{action: ActionStudio, wantArgs: []string{"prisma", "studio"}},
		{action: ActionSeed, wantArgs: []string{"prisma", "db", "seed"}},
		{action: ActionReset, wantArgs: []string{"prisma", "migrate", "reset"}},
	}

	for _, tt := range tests {
		t.Run(string(tt.action), func(t *testing.T) {
			cmd, err := ResolveCommand(pCtx, tt.action, nil)
			if err != nil {
				t.Fatalf("unexpected error for %s: %v", tt.action, err)
			}
			if cmd.Executable != "npx" {
				t.Errorf("expected npx, got %s", cmd.Executable)
			}
			if strings.Join(cmd.Args, " ") != strings.Join(tt.wantArgs, " ") {
				t.Errorf("args = %v, want %v", cmd.Args, tt.wantArgs)
			}
		})
	}
}

func TestResolveCommand_MongooseErrorWithoutScript(t *testing.T) {
	pCtx := &ProjectDBContext{
		ProjectDir:     "/test/mern",
		Layout:         "standalone",
		PackageManager: "npm",
		Provider:       "mongodb",
		ORM:            "mongoose",
		RootScripts:    map[string]string{},
	}

	_, err := ResolveCommand(pCtx, ActionMigrate, nil)
	if err == nil {
		t.Errorf("expected error for mongoose without script, got nil")
	}
}

type mockExecutor struct {
	executedCmd *ResolvedCommand
}

func (m *mockExecutor) Execute(ctx context.Context, cmd ResolvedCommand, stdin io.Reader, stdout, stderr io.Writer) error {
	m.executedCmd = &cmd
	return nil
}

func TestRunner_DryRun(t *testing.T) {
	tmpDir := t.TempDir()

	configJSON := `{
		"architecture": {"layout": "standalone", "packageManager": "bun"},
		"stack": {"database": {"provider": "sqlite", "orm": "drizzle"}}
	}`
	_ = os.WriteFile(filepath.Join(tmpDir, "koko.config.json"), []byte(configJSON), 0644)

	var stdoutBuf bytes.Buffer
	mockExec := &mockExecutor{}
	runner := &Runner{
		Executor: mockExec,
	}

	err := runner.Run(context.Background(), Options{
		Dir:    tmpDir,
		Action: ActionPush,
		DryRun: true,
		Stdout: &stdoutBuf,
	})
	if err != nil {
		t.Fatalf("unexpected dry-run error: %v", err)
	}

	output := stdoutBuf.String()
	if !strings.Contains(output, "Dry Run Analysis") {
		t.Errorf("expected dry run header in output: %s", output)
	}
	if !strings.Contains(output, "bunx drizzle-kit push") {
		t.Errorf("expected resolved bunx drizzle-kit push in output: %s", output)
	}
	if mockExec.executedCmd != nil {
		t.Errorf("expected no execution during dry-run")
	}
}

func TestPrintDBStackOverview(t *testing.T) {
	tmpDir := t.TempDir()

	configJSON := `{
		"architecture": {"layout": "monorepo", "packageManager": "pnpm"},
		"stack": {"database": {"provider": "postgres", "orm": "prisma"}}
	}`
	_ = os.WriteFile(filepath.Join(tmpDir, "koko.config.json"), []byte(configJSON), 0644)

	var stdoutBuf bytes.Buffer
	err := PrintDBStackOverview(&stdoutBuf, tmpDir)
	if err != nil {
		t.Fatalf("unexpected overview error: %v", err)
	}

	out := stdoutBuf.String()
	if !strings.Contains(out, "Database Stack Overview") {
		t.Errorf("expected title in overview: %s", out)
	}
	if !strings.Contains(out, "postgres") || !strings.Contains(out, "prisma") {
		t.Errorf("expected postgres and prisma in overview: %s", out)
	}
	if !strings.Contains(out, "migrate") || !strings.Contains(out, "studio") {
		t.Errorf("expected subcommands listed: %s", out)
	}
}

func TestResolveCommand_MonorepoPackagesDBPriority(t *testing.T) {
	tmpDir := t.TempDir()

	// Monorepo layout with packages/db
	packagesDbDir := filepath.Join(tmpDir, "packages", "db")
	if err := os.MkdirAll(packagesDbDir, 0755); err != nil {
		t.Fatalf("failed creating packages/db: %v", err)
	}

	rootPkg := `{"name": "root-repo", "scripts": {}}`
	_ = os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(rootPkg), 0644)

	dbPkg := `{"name": "@repo/db", "scripts": {"db:migrate": "drizzle-kit migrate", "db:push": "drizzle-kit push"}}`
	_ = os.WriteFile(filepath.Join(packagesDbDir, "package.json"), []byte(dbPkg), 0644)

	configJSON := `{
		"architecture": {"layout": "monorepo", "packageManager": "pnpm"},
		"stack": {"database": {"provider": "postgres", "orm": "drizzle"}}
	}`
	_ = os.WriteFile(filepath.Join(tmpDir, "koko.config.json"), []byte(configJSON), 0644)

	pCtx, err := DetectProjectDBContext(tmpDir)
	if err != nil {
		t.Fatalf("unexpected detect error: %v", err)
	}
	if !pCtx.HasDBPackage {
		t.Fatalf("expected HasDBPackage = true")
	}

	resCmd, err := ResolveCommand(pCtx, ActionMigrate, nil)
	if err != nil {
		t.Fatalf("unexpected resolve error: %v", err)
	}

	if resCmd.Executable != "pnpm" {
		t.Errorf("exe = %s, want pnpm", resCmd.Executable)
	}
	// Debería usar el filtro de monorepo pnpm --filter @*/db run db:migrate
	expectedArgs := []string{"--filter", "@*/db", "run", "db:migrate"}
	if strings.Join(resCmd.Args, " ") != strings.Join(expectedArgs, " ") {
		t.Errorf("args = %v, want %v", resCmd.Args, expectedArgs)
	}
}

func TestResolveCommand_Security_NoShellCommandInjection(t *testing.T) {
	// Verificar que argumentos potencialmente maliciosos permanezcan como strings estructurados en cmd.Args
	pCtx := &ProjectDBContext{
		ProjectDir:     "/test/standalone",
		Layout:         "standalone",
		PackageManager: "pnpm",
		Provider:       "postgres",
		ORM:            "drizzle",
		RootScripts:    map[string]string{},
	}

	maliciousArgs := []string{"; rm -rf /", "&& cat /etc/passwd", "| echo hacked", "`whoami`"}
	cmd, err := ResolveCommand(pCtx, ActionPush, maliciousArgs)
	if err != nil {
		t.Fatalf("unexpected error resolving command: %v", err)
	}

	for _, malArg := range maliciousArgs {
		found := false
		for _, arg := range cmd.Args {
			if arg == malArg {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected malicious argument %q to be preserved literally as structured arg", malArg)
		}
	}
}

