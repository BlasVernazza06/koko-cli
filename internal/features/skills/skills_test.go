package skills

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
)

func TestDetectRunner(t *testing.T) {
	tests := []struct {
		name         string
		available    map[string]string
		wantName     string
		wantExe      string
		wantPrefix   []string
		wantErrIs    error
	}{
		{
			name: "npx is detected first",
			available: map[string]string{
				"npx":  "/usr/local/bin/npx",
				"pnpm": "/usr/local/bin/pnpm",
				"bun":  "/usr/local/bin/bun",
			},
			wantName:   "npx",
			wantExe:    "/usr/local/bin/npx",
			wantPrefix: []string{"autoskills"},
		},
		{
			name: "bunx is detected if npx is missing",
			available: map[string]string{
				"bunx": "/usr/local/bin/bunx",
				"npm":  "/usr/local/bin/npm",
			},
			wantName:   "bunx",
			wantExe:    "/usr/local/bin/bunx",
			wantPrefix: []string{"autoskills"},
		},
		{
			name: "pnpm dlx is used if npx and bunx are missing",
			available: map[string]string{
				"pnpm": "/usr/local/bin/pnpm",
				"npm":  "/usr/local/bin/npm",
			},
			wantName:   "pnpm",
			wantExe:    "/usr/local/bin/pnpm",
			wantPrefix: []string{"dlx", "autoskills"},
		},
		{
			name: "bun x is used if bun is available",
			available: map[string]string{
				"bun": "/usr/local/bin/bun",
			},
			wantName:   "bun",
			wantExe:    "/usr/local/bin/bun",
			wantPrefix: []string{"x", "autoskills"},
		},
		{
			name: "npm exec fallback is used",
			available: map[string]string{
				"npm": "/usr/local/bin/npm",
			},
			wantName:   "npm",
			wantExe:    "/usr/local/bin/npm",
			wantPrefix: []string{"exec", "autoskills"},
		},
		{
			name:      "no runner available returns ErrNoRunnerFound",
			available: map[string]string{},
			wantErrIs: ErrNoRunnerFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				LookPath: func(file string) (string, error) {
					if path, ok := tt.available[file]; ok {
						return path, nil
					}
					return "", errors.New("not found")
				},
				Exec: defaultCommandExecutor,
			}

			runner, err := m.DetectRunner()
			if tt.wantErrIs != nil {
				if !errors.Is(err, tt.wantErrIs) {
					t.Fatalf("expected error %v, got %v", tt.wantErrIs, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if runner.Name != tt.wantName {
				t.Errorf("runner.Name = %q, want %q", runner.Name, tt.wantName)
			}
			if runner.Executable != tt.wantExe {
				t.Errorf("runner.Executable = %q, want %q", runner.Executable, tt.wantExe)
			}
			if !reflect.DeepEqual(runner.PrefixArgs, tt.wantPrefix) {
				t.Errorf("runner.PrefixArgs = %v, want %v", runner.PrefixArgs, tt.wantPrefix)
			}
		})
	}
}

func TestBuildAutoskillsArgs(t *testing.T) {
	tests := []struct {
		name string
		opts Options
		want []string
	}{
		{
			name: "empty options",
			opts: Options{},
			want: nil,
		},
		{
			name: "yes flag",
			opts: Options{Yes: true},
			want: []string{"--yes"},
		},
		{
			name: "yes flag not duplicated if already in args",
			opts: Options{Yes: true, Args: []string{"--yes", "ui"}},
			want: []string{"--yes", "ui"},
		},
		{
			name: "short -y flag prevents adding --yes",
			opts: Options{Yes: true, Args: []string{"-y", "ui"}},
			want: []string{"-y", "ui"},
		},
		{
			name: "dry run flag",
			opts: Options{DryRun: true},
			want: []string{"--dry-run"},
		},
		{
			name: "dry run not duplicated",
			opts: Options{DryRun: true, Args: []string{"--dry-run"}},
			want: []string{"--dry-run"},
		},
		{
			name: "agent flag",
			opts: Options{Agent: "cursor"},
			want: []string{"--agent", "cursor"},
		},
		{
			name: "agent not duplicated if already in args with equal",
			opts: Options{Agent: "cursor", Args: []string{"--agent=cursor"}},
			want: []string{"--agent=cursor"},
		},
		{
			name: "all flags combined with subcommands",
			opts: Options{
				Yes:    true,
				DryRun: true,
				Agent:  "claude",
				Args:   []string{"ui"},
			},
			want: []string{"--yes", "--dry-run", "--agent", "claude", "ui"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildAutoskillsArgs(tt.opts)
			if len(got) == 0 && len(tt.want) == 0 {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildAutoskillsArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRun(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("fails when directory does not exist", func(t *testing.T) {
		m := NewManager()
		err := m.Run(context.Background(), Options{
			Dir: filepath.Join(tmpDir, "non_existent_dir_12345"),
		})
		if !errors.Is(err, ErrDirectoryNotFound) {
			t.Fatalf("expected ErrDirectoryNotFound, got %v", err)
		}
	})

	t.Run("fails when target is a file instead of directory", func(t *testing.T) {
		filePath := filepath.Join(tmpDir, "somefile.txt")
		if err := os.WriteFile(filePath, []byte("content"), 0644); err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}

		m := NewManager()
		err := m.Run(context.Background(), Options{
			Dir: filePath,
		})
		if !errors.Is(err, ErrDirectoryNotADir) {
			t.Fatalf("expected ErrDirectoryNotADir, got %v", err)
		}
	})

	t.Run("fails when no runner is found in system", func(t *testing.T) {
		m := &Manager{
			LookPath: func(file string) (string, error) {
				return "", errors.New("not found")
			},
			Exec: func(ctx context.Context, dir string, stdin io.Reader, stdout, stderr io.Writer, name string, args ...string) error {
				return nil
			},
		}

		err := m.Run(context.Background(), Options{
			Dir: tmpDir,
		})
		if !errors.Is(err, ErrNoRunnerFound) {
			t.Fatalf("expected ErrNoRunnerFound, got %v", err)
		}
	})

	t.Run("executes runner with correct directory, stdio, and arguments", func(t *testing.T) {
		var (
			executedDir  string
			executedName string
			executedArgs []string
		)

		var stdoutBuf bytes.Buffer
		var stderrBuf bytes.Buffer
		stdinBuf := strings.NewReader("input")

		m := &Manager{
			LookPath: func(file string) (string, error) {
				if file == "npx" {
					return "/path/to/npx", nil
				}
				return "", errors.New("not found")
			},
			Exec: func(ctx context.Context, dir string, stdin io.Reader, stdout, stderr io.Writer, name string, args ...string) error {
				executedDir = dir
				executedName = name
				executedArgs = args
				return nil
			},
		}

		err := m.Run(context.Background(), Options{
			Dir:    tmpDir,
			Yes:    true,
			DryRun: true,
			Agent:  "cursor",
			Args:   []string{"ui"},
			Stdin:  stdinBuf,
			Stdout: &stdoutBuf,
			Stderr: &stderrBuf,
		})
		if err != nil {
			t.Fatalf("expected success, got %v", err)
		}

		if executedDir != tmpDir {
			t.Errorf("executedDir = %q, want %q", executedDir, tmpDir)
		}
		if executedName != "/path/to/npx" {
			t.Errorf("executedName = %q, want %q", executedName, "/path/to/npx")
		}

		expectedArgs := []string{"autoskills", "--yes", "--dry-run", "--agent", "cursor", "ui"}
		if !reflect.DeepEqual(executedArgs, expectedArgs) {
			t.Errorf("executedArgs = %v, want %v", executedArgs, expectedArgs)
		}
	})

	t.Run("wraps command execution error properly", func(t *testing.T) {
		expectedErr := errors.New("command terminated with exit status 1")
		m := &Manager{
			LookPath: func(file string) (string, error) {
				return "/path/to/npx", nil
			},
			Exec: func(ctx context.Context, dir string, stdin io.Reader, stdout, stderr io.Writer, name string, args ...string) error {
				return expectedErr
			},
		}

		err := m.Run(context.Background(), Options{
			Dir: tmpDir,
		})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Fatalf("expected wrapped error %v, got %v", expectedErr, err)
		}
	})
}
