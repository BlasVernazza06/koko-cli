---
id: go-cli-architecture
title: Go CLI & Cobra Architecture
category: Go
description: Production-grade CLI development with Cobra, flags, command trees, in-memory VFS, and TUI integration.
tags: [go, cli, cobra, tui, architecture]
globs: ["**/*.go", "cmd/**/*.go"]
---

# Go CLI & Cobra Architecture Guidelines

Architecting high-performance CLI utilities in Go with `spf13/cobra` requires strict flag hygiene, sub-10ms startup latency, and atomic filesystem safety.

## Standard Command Hierarchy Reference

| Component | Layer | Purpose | Key Interface / Pattern |
| :--- | :--- | :--- | :--- |
| `cmd/root.go` | CLI Root | Root command definition, global flags (`--dir`, `--verbose`), version | `cobra.Command.Execute()` |
| `cmd/<command>.go` | CLI Branch | Subcommands (`init.go`, `add.go`, `doctor.go`, `skills.go`). Parses args | `RunE: func(cmd, args) error` |
| `internal/core/vfs/` | Buffer | In-memory Virtual File System (`sync.RWMutex`) preventing dirty disk states | `vfs.New()`, `vfs.WriteFile()` |
| `internal/engine/` | Business Core | Core scaffolding, dependency injection, and compatibility matrix | Stateless pure engine functions |
| `internal/tui/` | View Layer | Bubble Tea models, Lip Gloss rendering, keyboard event loops | `tea.Model`, `tea.Program` |

## CLI Flag Conventions & Exit Codes

| Convention | Rule | Bad Example | Good Example |
| :--- | :--- | :--- | :--- |
| **Short/Long Flags** | Provide single-letter shortcuts for common operations | `--directory` only | `-d, --dir` |
| **Affirmative Bypass** | Provide `-y, --yes` to skip prompts in non-interactive CI | Stalling pipeline on stdin | `if opts.Yes { proceed() }` |
| **Dry Run Mode** | Mutating commands must support `--dry-run` | Writing files without preview | Print actions, skip `WriteTree()` |
| **Exit Code 0** | Clean execution without errors | `os.Exit(1)` on help/version | `return nil` |
| **Exit Code 1+** | Handled domain errors with actionable messages | Raw panic stack trace | Clean error to `os.Stderr` |

## Cold-Start & Performance Rules

- **Zero Heavy Work in `init()`**: Do not load network sockets, parse heavy configurations, or spawn subshells in `init()`. Defer initialization until `PreRunE` or `RunE`.
- **Target Startup Latency**: A CLI command checking `--help` or `--version` must return in under 10ms.
- **Lazy Loading**: Use lazily initialized package managers or singletons when external resources (e.g., autoskills runners or Docker daemons) might not be needed by the executed subcommand.

## Anti-Pattern vs Recommended Solution

### ❌ Anti-Pattern: Direct Filesystem Mutations During Execution
Writing directly to disk during generation means an error halfway leaves partial, corrupted directories on the user's computer.

```go
// BAD: Modifying disk directly step-by-step
func generateProject(name string) error {
    if err := os.Mkdir(name, 0755); err != nil {
        return err
    }
    // Writing files one by one to disk
    os.WriteFile(filepath.Join(name, "package.json"), pkgData, 0644)
    
    // If validation fails halfway through, the user is left with a corrupted directory!
    if err := validateStack(); err != nil {
        return err // Dirty corrupted state on disk!
    }
    return nil
}
```

### ✅ Recommended: In-Memory VFS with Atomic Safe Flush

```go
// GOOD: Generate and validate completely in memory first; flush to disk only upon 100% success
func GenerateProject(name string, cfg Config) error {
    memFS := vfs.New()

    // 1. Render all files in memory (Zero disk side-effects)
    memFS.WriteFile("package.json", pkgData)
    memFS.WriteFile("koko.config.json", configData)
    memFS.WriteFile(".env", envData)

    // 2. Validate consistency in memory
    if err := ValidateVFS(memFS); err != nil {
        return fmt.Errorf("validation error in memory: %w", err) // Disk remains 100% pristine!
    }

    // 3. Atomic flush to disk
    return vfs.WriteTree(memFS, name)
}
```

## Anti-Pattern: Blocking Stdin in Non-TTY Environments

### ❌ Anti-Pattern: Prompting without TTY Detection
Spawning interactive prompts when executed in CI/CD pipelines hangs automated builds indefinitely.

```go
// BAD: Attempting to prompt stdin in automated environments
func promptUserChoice() string {
    var input string
    fmt.Scanln(&input) // Blocks forever on headless GitHub Actions or scripts!
    return input
}
```

### ✅ Recommended: Non-Interactive Flag with Fallback

```go
// GOOD: Detect flags or headless environments and fallback gracefully
func ResolveUserSelection(cmd *cobra.Command, opts Options) (string, error) {
    if opts.Yes || !isatty.IsTerminal(os.Stdin.Fd()) {
        return DefaultSelection, nil
    }
    return RunInteractiveBubbleTea(cmd.Context())
}
```
