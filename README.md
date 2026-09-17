<div align="center">

# <img src="assets/koko.png" alt="Koko CLI Logo" width="38" style="vertical-align: middle;" /> Koko CLI

### *Grab your stack, structure your project, and start building instantly.*

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

</div>

---

**Koko CLI** is an ultra-fast interactive terminal project initializer written in Go. Forget about manually setting up complex TypeScript configurations, monorepo workspaces, databases, or Docker containers. Koko lets you pick your ideal production recipe or design a custom stack step-by-step, generating ready-to-run projects in milliseconds.

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

### Instant Run (Zero Installation)
Bootstrap a production-ready project immediately with your favorite package manager:

```bash
# Node.js (NPX)
npx koko-app init

# PNPM
pnpm dlx koko-app init

# Bun
bunx koko-app init
```

### Native Go Installation
```bash
go install github.com/BlasVernazza06/koko-cli@latest
```

### Global CLI Installation
```bash
npm install -g koko-app
# or
pnpm add -g koko-app
# or
bun add -g koko-app
```

---

## 🔮 What's New in v2.0.0

* ⚡ **100% Go-Template Engine**: All codebase templates, configurations, environment variables, and SDK clients are pre-compiled and interpolated directly at generation time. Zero post-install tweaks required—run `pnpm install` and your monorepo is ready to develop.
* 🧩 **Smart `koko add` Command**: Expand your monorepo anytime without starting from scratch.
  * **Fresh Projects (No custom code changes):** Performs a full scaffold of the selected addon (e.g. `koko add stripe`, `koko add clerk`, `koko add zod`), injecting client libraries, database schemas, and example API routes.
  * **Modified Projects (Code changes detected):** Safely injects dependencies and `.env.example` keys without overwriting your custom logic.
  * **New Apps / Workspaces:** If you started with only a frontend or backend, `koko add backend` or `koko add mobile` creates the new workspace folder (`apps/api`, `apps/mobile`) and links it to the monorepo automatically.
* 🛡️ **Workspace Integrity Engine**: Real-time cross-validation prevents invalid architectural pairings and guarantees 100% valid `workspace:*` links.

---

## 🛠️ Supported Technology Stack

Koko supports a diverse set of production-ready technologies:

### 📦 1. Pre-configured Production Recipes
* ⚡ **SaaS Starter:** Next.js + Drizzle ORM + Better-Auth + Docker Compose + Stripe.
* 💻 **MERN Stack:** React (Vite) + Express + MongoDB.
* 🚀 **PERN Stack:** React (Vite) + Express + PostgreSQL.
* 🐍 **FastAPI + React:** Python FastAPI + React (Vite) SPA.
* ☕ **Enterprise NestJS:** NestJS API + Next.js + Prisma + PostgreSQL.
* 📱 **Mobile Expo:** React Native Expo + Express Backend + Shared DB.

### ⚙️ 2. Step-by-Step Modular Stacks

| Layer | Supported Technologies |
| :--- | :--- |
| **Frontend Framework** | `Next.js` (App Router) • `React (Vite)` • `Astro` • `Nuxt (Vue)` • `Svelte` • `None` |
| **Backend Runtime** | `Node.js Express` • `NestJS` • `Hono` • `Go Chi Router` • `Python FastAPI` • `None` |
| **Database Server** | `PostgreSQL` • `MySQL` • `MongoDB` • `SQLite` • `None` |
| **ORM / Query Builder** | `Drizzle ORM` • `Prisma` • `Mongoose` • `GORM (Go)` • `SQLAlchemy / SQLModel (Python)` • `None` |
| **Authentication** | `Better-Auth` • `Clerk` • `NextAuth` • `None` |
| **API Layer** | `tRPC` • `oRPC` • `REST` |
| **Package Managers** | `pnpm` • `bun` • `npm` • `Go Modules` • `pip` |
| **Addons & Tooling** | `Shadcn UI` • `Stripe` • `Polar` • `Zod` • `Docker Compose` • `GitHub Actions CI` |

---

## 🧩 The `koko add` Workflow

Forgot to select an addon during initialization? Add it in seconds:

```bash
# Add payments & billing
koko add stripe
koko add polar

# Add authentication
koko add clerk
koko add better-auth

# Add UI or schema validators
koko add shadcn
koko add zod

# Add a new workspace application to an existing monorepo
koko add backend    # Scaffolds apps/api and links workspace
koko add mobile     # Scaffolds apps/mobile (Expo) and links workspace
```

---

## 📁 Generated Monorepo Layout

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

Koko stores project configuration in a root manifest to power Day-2 operations (`koko add`, `koko doctor`, etc.):

```json
{
  "$schema": "https://koko-cli.dev/schema.json",
  "project": {
    "name": "my-super-app",
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

Contributions are welcome! Check out our [Contributing Guide](CONTRIBUTING.md) to get started.

## 📄 License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
