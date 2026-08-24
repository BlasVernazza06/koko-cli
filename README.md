<div align="center">

# 🦾 Koko CLI

### *Grab your stack, structure your project, and start building instantly.*

[![npm version](https://img.shields.io/npm/v/koko-cli?color=33cd5f&style=flat-flat)](https://www.npmjs.com/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

</div>

---

**Koko CLI** (guided by your faithful scaffolding assistant 🦾) is an ultra-fast interactive terminal initializer written in Go. Forget about manually setting up complex TypeScript templates, monorepo workspaces, databases, or Docker containers. Koko lets you pick your ideal production recipe and structures your project ready for development in seconds, 100% offline.

```text
     .---.   .---.
    /  _  \_/  _  \       _  ______  _  ______
   |  (o)     (o)  |     | |/ / __ \| |/ / __ \
   |     (..)      |     | ' / /  | | ' / /  | |
    \   (____)    /      | . \ \__| | . \ \__| |
     '-----------'       |_|\_\____/|_|\_\____/  v1.0.0
```

---

## ✨ Key Features

*   🎨 **Premium Visual Experience**: Interactive, animated TUI built with Bubble Tea and Lipgloss.
*   🚀 **Instant Scaffolding**: Physical disk unpacking and structuring in milliseconds.
*   📦 **Production Recipes**: Bootstrapping with ready-to-use, optimized setups (e.g. SaaS Starter with auth and DB).
*   🐳 **Monorepo & Infrastructure**: Automatic configuration of pnpm/npm workspaces, linters, TypeScript, and local Docker Compose containers.

---

## 🛠️ Supported Technology Stack (v1.0.0)

| Frontend 💻 | Backend 🔌 | Infrastructure 🐳 | Database & ORM 🗄️ |
| :--- | :--- | :--- | :--- |
| **Next.js** (App Router, TS) | **Node.js Express** (TS) | **Docker Compose** | **PostgreSQL** (Drizzle / Prisma) |
| **React** (Vite SPA, TS) | **Go Fiber** (REST API) | **pnpm Workspaces** | **MySQL** (Drizzle / Prisma / SQLx) |
| **Vue.js** (Vite, TS) | **Hono** (Node.js Server) | **ESLint & Prettier** | **MongoDB** (Mongoose) |

---

## 🚀 Quick Start

### Run Directly (Without Installing)
The fastest way to start a new project is with `npx`:
```bash
npx koko-cli init
```

### 🛠️ Local Development & Build (Go)
If running or modifying the CLI from source:
```bash
# Build binary
go build -o koko.exe main.go

# Run interactive assistant
./koko.exe init
```

---

## 💻 Terminal Flow

When started, **Koko** guides you with interactive progress indicators:

```bash
  ◇ Scaffolding directories and copying templates...  ✓
  ◇ Generating Docker and Database configuration...   ✓
  ◇ Initializing Git repository...                    ✓
  ◇ Creating koko.config.json manifest...             ✓

✨ Project 'my-app' created successfully in 2.45s!
```

---

## 📁 Generated Repository Structure

Depending on your choices, Koko structures a unified modular monorepo architecture:

```text
my-app/
├── apps/
│   ├── frontend/         # Frontend application (e.g., Next.js)
│   └── backend/          # API server (e.g., Express / Go Fiber)
├── packages/
│   └── db/               # Schemas and connection client (Prisma / Drizzle)
├── docker-compose.yml    # Local PostgreSQL / MySQL container
├── koko.config.json      # Koko configuration manifest
├── package.json          # Root package with configured workspaces
└── README.md
```

---

## ⚙️ Configuration Manifest (koko.config.json)

When initializing a project, a `koko.config.json` file is generated at the root, acting as the single source of truth for the stack and enabling "Day 2" additions (e.g., via `koko add`):

```json
{
  "$schema": "https://koko-cli.dev/schema.json",
  "project": {
    "name": "my-super-app",
    "cliVersion": "v1.0.0",
    "createdAt": "2026-08-20T20:50:00Z"
  },
  "architecture": {
    "layout": "monorepo",
    "packageManager": "pnpm"
  },
  "stack": {
    "frontend": {
      "framework": "next",
      "language": "typescript",
      "styling": "tailwindcss"
    },
    "backend": {
      "framework": "express",
      "language": "typescript"
    },
    "database": {
      "provider": "postgres",
      "orm": "prisma"
    }
  },
  "features": {
    "auth": {
      "provider": "better-auth",
      "status": "installed"
    },
    "infrastructure": {
      "dockerCompose": true,
      "ciCd": "github-actions"
    }
  }
}
```

---

## 🤝 Contributing

Contributions are welcome! Check our [Contributing Guide](CONTRIBUTING.md) to learn how to clone the repository, set up your local environment, and submit Pull Requests.

---

## 📄 License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

