# 💻 `cmd` — CLI Command Layer & Entrypoints

This directory contains the user-facing CLI command layer for **Koko CLI**, built with [spf13/cobra](https://github.com/spf13/cobra) and [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea).

---

## 🗂️ Architecture & Subcommands

| File / Directory | Purpose | Responsibility |
| :--- | :--- | :--- |
| `root.go` | Root Cobra Command (`koko`) | Global flags, version output, help banners, and command registration. |
| `init.go` | Initialization Command (`koko init`) | Parses flags (`--recipe`, `--pm`, `--yes`, etc.), executes non-interactive flows or hands over execution to the interactive TUI. |
| `doctor.go` | Diagnostics Command (`koko doctor`) | Inspects local developer dependencies (Node.js, Go, Docker, Bun, PNPM) and workspace integrity. |
| `tui.go` | Terminal UI Dispatcher | Initializes and runs the Bubble Tea interactive loop when commands run in wizard mode. |
| `views/` | TUI Screens & Lipgloss Styles | Modular views, banners, spinners, and form selectors for the interactive experience. |

---

## 🚀 Key Guidelines
* **POSIX Compliance**: Ensure all flags follow standard naming conventions (kebab-case, single-letter shorthands where appropriate).
* **Decoupled Business Logic**: Command handlers should parse arguments/flags and delegate heavy scaffolding logic to `internal/scaffold` and `internal/vfs`.
