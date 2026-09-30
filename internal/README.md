# ⚙️ `internal` — Core Scaffolding & Validation Engine

This directory contains the private core packages of **Koko CLI**. Code here cannot be imported by external Go modules.

---

## 🗂️ Package Overview

```text
internal/
├── catalog/        # Registered recipes, presets, and stack definitions
├── compatibility/  # Compatibility engine and invalid combination validator
├── config/         # Persistent configuration and preferences
├── doctor/         # Health-check diagnostics runners for external toolchains
├── errors/         # Custom domain error types and user-friendly formatting
├── scaffold/       # In-memory template rendering, file generation, and recipe execution
├── types/          # Shared domain structs, enums, and manifest schemas
├── ui/             # Lipgloss themes, icons, terminal styling helpers
├── validator/      # Monorepo and workspace link integrity checkers
└── vfs/            # Virtual File System (In-memory tree before disk write)
```

---

## 🛡️ Core Architectural Principles
1. **In-Memory VFS First**: Templates and files are staged in memory in `internal/vfs` before performing any disk I/O. This ensures atomic rollbacks if an error occurs.
2. **Compatibility Enforcement**: `internal/compatibility` validates frontend/backend/ORM combinations before any file is generated.
3. **Deterministic Output**: Given the same recipe and options, the scaffolding engine produces identical, clean workspace structures.
