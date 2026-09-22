<div align="center">

# <img src="assets/koko.png" alt="Koko CLI Logo" width="42" style="vertical-align: middle;" /> Koko CLI

### *Grab your stack, structure your project, and start building instantly.*

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://golang.org)
[![NPM Version](https://img.shields.io/npm/v/koko-app?style=flat-square&color=CB3837&logo=npm)](https://www.npmjs.com/package/koko-app)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](LICENSE)
[![Turborepo](https://img.shields.io/badge/Monorepo-Turborepo-EF4444?style=flat-square&logo=turborepo&logoColor=white)](https://turbo.build)

<br />

<p align="center">
  <img src="assets/demo.gif" alt="Koko CLI Interactive Demo" width="92%" style="border-radius: 8px; box-shadow: 0 8px 30px rgba(0,0,0,0.25);" />
</p>

</div>

---

**Koko CLI** is an ultra-fast interactive terminal project initializer written in Go. Forget about manually wiring TypeScript monorepo workspaces, database schemas, Docker containers, authentication providers, and client SDKs. Koko generates clean, fully-typed, production-ready fullstack architectures in milliseconds.

```text
     .---.   .---.
    /  _  \_/  _  \       _  ______  _  ______
   |  (o)     (o)  |     | |/ / __ \| |/ / __ \
   |     (..)      |     | ' / /  | | ' / /  | |
    \   (____)    /      | . \ \__| | . \ \__| |
     '-----------'       |_|\_\____/|_|\_\____/  v2.0.0
```

---

## ⚡ Quick Start

Bootstrap a production-ready project immediately with zero prior setup:

### Instant Execution

```bash
# Node.js (NPX)
npx koko-app init

# PNPM
pnpm dlx koko-app init

# Bun
bunx koko-app init
```

### Global / Native Installation

```bash
# Go Native Binary
go install github.com/BlasVernazza06/koko-cli@latest

# Global NPM / PNPM / Bun
npm install -g koko-app
# or
pnpm add -g koko-app
# or
bun add -g koko-app
```

Then run anywhere:
```bash
koko init
```

---

## ✨ Key Features

* ⚡ **Ultra-Fast In-Memory VFS**: All project files, configurations, and templates are rendered in memory before touching the disk, creating clean projects in under a second.
* 🛡️ **Workspace Integrity Engine**: Built-in real-time validation guarantees valid `workspace:*` links, preventing incompatible architectural combinations.
* 🎨 **Interactive TUI**: Intuitive terminal UI built with Charm Bubbletea featuring smart recipes, step-by-step custom wizard, and live progress indicators.
* 🧩 **Modular Day-2 Addons (`koko add`)**: Expand existing projects seamlessly with authentication, payments, database models, or new workspace apps (`apps/api`, `apps/mobile`) without breaking custom logic.
* 🩺 **Health Check & Diagnostics (`koko doctor`)**: Inspect dependencies, catalog updates, and monorepo workspace alignment in a single command.

---

## 🛠️ Supported Technology Stack

### 📦 1. Production-Ready Recipes
* ⚡ **SaaS Starter:** Next.js App Router + Drizzle ORM + PostgreSQL + Better-Auth + Docker Compose + Stripe.
* 🏢 **Enterprise NestJS:** Next.js + NestJS API + PostgreSQL (Prisma) + Better-Auth + Stripe.
* 🚀 **PERN Stack:** React (Vite) + Node.js Express + PostgreSQL + Prisma ORM.
* 💻 **MERN Stack:** React (Vite) + Node.js Express + MongoDB + Mongoose.
* 🐍 **FastAPI + React:** React SPA + Python FastAPI + PostgreSQL + Docker.
* 📱 **Mobile Expo:** React Native Expo + Express Backend + Shared Database.

### ⚙️ 2. Modular Stack Matrix

| Layer | Supported Technologies |
| :--- | :--- |
| **Frontend** | `Next.js` (App Router) • `React (Vite)` • `Astro` • `Nuxt (Vue)` • `Svelte` • `None` |
| **Backend** | `Node.js Express` • `NestJS` • `Hono` • `Go Chi Router` • `Python FastAPI` • `None` |
| **Database** | `PostgreSQL` • `MySQL` • `MongoDB` • `SQLite` • `None` |
| **ORM / Client** | `Drizzle ORM` • `Prisma` • `Mongoose` • `GORM (Go)` • `SQLAlchemy / SQLModel (Python)` • `None` |
| **Authentication** | `Better-Auth` • `Clerk` • `NextAuth.js` • `None` |
| **API Layer** | `tRPC` • `oRPC` • `REST` |
| **Package Managers**| `pnpm` • `bun` • `npm` • `Go Modules` • `pip` |
| **Addons & Tooling**| `Shadcn UI` • `Stripe` • `Polar` • `Zod` • `Docker Compose` • `GitHub Actions CI` |

---

## 🧩 The `koko add` Workflow

Need to add new capabilities to an existing workspace? Use `koko add`:

```bash
# Add payment & subscription providers
koko add stripe
koko add polar

# Add authentication
koko add clerk
koko add better-auth

# Add UI components & validation
koko add shadcn
koko add zod

# Expand monorepo workspaces
koko add backend    # Scaffolds apps/api and links workspace packages
koko add mobile     # Scaffolds apps/mobile (Expo) and connects to shared DB
```

---

## 📁 Generated Monorepo Architecture

```text
my-koko-app/
├── apps/
│   ├── web/              # Frontend application (Next.js, React, Astro, etc.)
│   └── api/              # Backend server (Express, Hono, Go Chi, FastAPI, NestJS)
├── packages/
│   ├── auth/             # Shared authentication configuration
│   ├── db/               # Shared Database client & ORM models (Prisma / Drizzle)
│   ├── ui/               # Shared component library (Shadcn UI)
│   ├── typescript-config/# Shared TypeScript configs across monorepo
│   └── eslint-config/    # Shared ESLint rules
├── docker-compose.yml    # Local Database container
├── koko.config.json      # Workspace configuration manifest
├── package.json          # Root workspaces configuration
├── turbo.json            # Turborepo pipeline orchestration
└── README.md
```

---

## ⚙️ Configuration Manifest (`koko.config.json`)

Koko records your stack specification in a root manifest to automate Day-2 maintenance, tooling, and doctor validations:

```json
{
  "$schema": "https://koko-cli.dev/schema.json",
  "project": {
    "name": "my-modern-stack",
    "cliVersion": "v2.0.0",
    "createdAt": "2026-09-17T20:50:00Z"
  },
  "architecture": {
    "layout": "monorepo",
    "packageManager": "pnpm"
  },
  "stack": {
    "frontend": {
      "framework": "nextjs",
      "language": "typescript"
    },
    "backend": {
      "framework": "go_chi",
      "language": "go"
    },
    "database": {
      "provider": "postgres",
      "orm": "drizzle"
    },
    "auth": "better-auth"
  },
  "features": {
    "infrastructure": {
      "dockerCompose": true,
      "ciCd": "github_actions"
    }
  }
}
```

---

## 🤝 Contributing

Contributions, issues, and feature requests are welcome! Feel free to check the [Contributing Guide](CONTRIBUTING.md).

## 📄 License

Distributed under the MIT License. See [LICENSE](LICENSE) for more information.

