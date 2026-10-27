package skills

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	// ErrNoRunnerFound is returned when neither npx, bunx, pnpm, bun, nor npm is available in PATH.
	ErrNoRunnerFound = errors.New("no Node.js package runner found (npx, pnpm, bun, or npm). Please ensure Node.js (with npx) or bun/pnpm is installed and available in your PATH")

	// ErrDirectoryNotFound is returned when the specified target directory does not exist.
	ErrDirectoryNotFound = errors.New("target directory does not exist")

	// ErrDirectoryNotADir is returned when the target path is not a directory.
	ErrDirectoryNotADir = errors.New("target path is not a directory")
)

// RunnerInfo contains metadata about the detected package runner executable.
type RunnerInfo struct {
	Name       string   // e.g. "npx", "bunx", "pnpm", "bun", "npm"
	Executable string   // Resolved path or binary name
	PrefixArgs []string // Arguments required to run autoskills (e.g. ["autoskills"] or ["dlx", "autoskills"])
}

// CommandExecutor abstracts command execution for testability and custom IO redirection.
type CommandExecutor func(ctx context.Context, dir string, stdin io.Reader, stdout, stderr io.Writer, name string, args ...string) error

// LookPathFunc abstracts binary path lookup in PATH.
type LookPathFunc func(file string) (string, error)

// Options specifies configuration for the skills execution.
type Options struct {
	Dir    string
	Yes    bool
	DryRun bool
	Agent  string
	Args   []string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// Manager coordinates runner detection and autoskills execution.
type Manager struct {
	LookPath LookPathFunc
	Exec     CommandExecutor
}

// defaultCommandExecutor executes the command using os/exec with terminal streams attached.
func defaultCommandExecutor(ctx context.Context, dir string, stdin io.Reader, stdout, stderr io.Writer, name string, args ...string) error {
	var cmd *exec.Cmd

	// On Windows, executing batch/cmd files (.cmd/.bat) directly via CreateProcess can fail
	// or require cmd.exe /c wrapping depending on environment and file path.
	if runtime.GOOS == "windows" && (strings.HasSuffix(strings.ToLower(name), ".cmd") || strings.HasSuffix(strings.ToLower(name), ".bat")) {
		cmdArgs := append([]string{"/c", name}, args...)
		cmd = exec.CommandContext(ctx, "cmd.exe", cmdArgs...)
	} else {
		cmd = exec.CommandContext(ctx, name, args...)
	}

	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd.Run()
}

// NewManager creates a Manager configured with default OS lookups and process execution.
func NewManager() *Manager {
	return &Manager{
		LookPath: exec.LookPath,
		Exec:     defaultCommandExecutor,
	}
}

// DetectRunner inspects the system to find an available Node/JS package runner in priority order.
func (m *Manager) DetectRunner() (*RunnerInfo, error) {
	candidates := []struct {
		name   string
		prefix []string
	}{
		{"npx", []string{"autoskills"}},
		{"bunx", []string{"autoskills"}},
		{"pnpm", []string{"dlx", "autoskills"}},
		{"bun", []string{"x", "autoskills"}},
		{"npm", []string{"exec", "autoskills"}},
	}

	for _, c := range candidates {
		if path, err := m.LookPath(c.name); err == nil && path != "" {
			return &RunnerInfo{
				Name:       c.name,
				Executable: path,
				PrefixArgs: c.prefix,
			}, nil
		}
	}

	return nil, ErrNoRunnerFound
}

// BuildAutoskillsArgs builds the list of CLI arguments to pass to autoskills,
// appending flags and avoiding duplicates if already provided in Args.
func BuildAutoskillsArgs(opts Options) []string {
	var result []string

	hasFlag := func(flags ...string) bool {
		for _, arg := range opts.Args {
			for _, f := range flags {
				if arg == f || strings.HasPrefix(arg, f+"=") {
					return true
				}
			}
		}
		return false
	}

	if opts.Yes && !hasFlag("--yes", "-y") {
		result = append(result, "--yes")
	}
	if opts.DryRun && !hasFlag("--dry-run") {
		result = append(result, "--dry-run")
	}
	if opts.Agent != "" && !hasFlag("--agent", "-a") {
		result = append(result, "--agent", opts.Agent)
	}

	result = append(result, opts.Args...)
	return result
}

// Run executes autoskills in the target directory using the detected runner.
func (m *Manager) Run(ctx context.Context, opts Options) error {
	targetDir := opts.Dir
	if targetDir == "" {
		targetDir = "."
	}
	cleanDir := filepath.Clean(targetDir)

	info, err := os.Stat(cleanDir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%w: %s", ErrDirectoryNotFound, cleanDir)
		}
		return fmt.Errorf("failed to check target directory '%s': %w", cleanDir, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%w: %s", ErrDirectoryNotADir, cleanDir)
	}

	runner, err := m.DetectRunner()
	if err != nil {
		return err
	}

	fullArgs := append([]string{}, runner.PrefixArgs...)
	fullArgs = append(fullArgs, BuildAutoskillsArgs(opts)...)

	stdin := opts.Stdin
	if stdin == nil {
		stdin = os.Stdin
	}
	stdout := opts.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = os.Stderr
	}

	if err := m.Exec(ctx, cleanDir, stdin, stdout, stderr, runner.Executable, fullArgs...); err != nil {
		return fmt.Errorf("autoskills execution failed: %w", err)
	}

	return nil
}

// DetectRunner is a package-level helper that invokes NewManager().DetectRunner().
func DetectRunner() (*RunnerInfo, error) {
	return NewManager().DetectRunner()
}

// Run is a package-level helper that invokes NewManager().Run().
func Run(ctx context.Context, opts Options) error {
	return NewManager().Run(ctx, opts)
}
