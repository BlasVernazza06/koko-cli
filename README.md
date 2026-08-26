<div align="center">

# 🦾 Koko CLI

### *Grab your stack, structure your project, and start building instantly.*

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

</div>

---

**Koko CLI** is an ultra-fast interactive terminal project initializer written in Go. Forget about manually setting up complex TypeScript configurations, monorepo workspaces, databases, or Docker containers. Koko lets you pick your ideal production recipe or design a custom stack step-by-step, generating your files in milliseconds.

```text
     .---.   .---.
    /  _  \_/  _  \       _  ______  _  ______
   |  (o)     (o)  |     | |/ / __ \| |/ / __ \
   |     (..)      |     | ' / /  | | ' / /  | |
    \   (____)    /      | . \ \__| | . \ \__| |
     '-----------'       |_|\_\____/|_|\_\____/  v2.0.0
```

---

## 📺 Demo walkthrough

Here is Koko CLI in action, generating a fully configured monorepo using the step-by-step TUI wizard:

<video src="demo.mp4" autoplay loop muted playsinline controls class="d-block rounded-2 border width-fit" style="max-height:640px; width: 100%;"></video>

> [!TIP]
> **Want to record your own demo?** Install [vhs](https://github.com/charmbracelet/vhs) and run `vhs demo.tape` from the root of this repository to generate a fresh `demo.mp4` automatically!

---

## ✨ Key Features

* 🎨 **Premium Interactive TUI**: Beautiful terminal interface built with Bubble Tea and Lipgloss featuring real-time spinner animations.
* 📦 **Flexible Init Modes**: Choose **Quick Setup** for pre-configured production recipes or **Manual Configuration** to build your stack piece-by-piece.
* 🛡️ **Cross-Validation Safety Rules**: Smart validation checks prevent generating incompatible tech combinations (e.g. blocking Go Fiber with Prisma ORM or Node-only drivers with Python) in real-time.
* 🐳 **Monorepo Architecture out-of-the-box**: Automatically setups workspaces, Shared TypeScript configs, ESLint/Prettier, local Docker Compose files, and CI/CD pipelines.

---

## 🛠️ Supported Technology Stack (v2.0.0)

Koko supports a diverse set of technologies, frameworks, and databases:

### 📦 1. Quick Setup Recipes
* ⚡ **SaaS Starter:** Next.js + Drizzle ORM + Better-Auth + Docker Compose + Stripe.
* 💻 **MERN Stack:** React (Vite) + Express + MongoDB.
* 🚀 **PERN Stack:** React (Vite) + Express + PostgreSQL.
* 🐍 **FastAPI + React:** Python FastAPI + React (Vite) SPA.

### ⚙️ 2. Manual Configuration Options

| Layer | Supported Technologies |
| :--- | :--- |
| **Frontend Framework** | `Next.js` (App Router) • `React (Vite)` • `Nuxt (Vue)` • `Svelte` • `None` |
| **Backend Runtime** | `Node.js Express` • `NestJS` • `Hono` • `Go Chi Router` • `Python FastAPI` • `None` |
| **Database Server** | `PostgreSQL` • `MySQL` • `MongoDB` • `SQLite` • `None` |
| **ORM / Query Builder** | `Drizzle ORM` • `Prisma` • `Mongoose` • `GORM (Go)` • `SQLAlchemy / SQLModel (Python)` • `None` |
| **Package Managers** | `pnpm` • `npm` • `bun` • `Go Modules` • `pip (Python)` |
| **Addons & Tooling** | `Docker Compose` • `GitHub Actions CI` • `None` |

---

## 🚀 Quick Start

### Run Directly via Node/NPX
You can run the CLI immediately without manual setup:
```bash
npx koko-cli init
```

### Download Native Binary (Go)
Download the pre-compiled executable matching your operating system directly from our **GitHub Releases** page and run it:
```bash
# Windows
koko.exe init

# macOS / Linux
chmod +x koko
./koko init
```

---

## 📁 Generated Repository Layout

Manual configurations structure a clean monorepo using pnpm/npm workspaces:

```text
my-koko-app/
├── apps/
│   ├── web/              # Frontend application (Next.js, React, etc.)
│   └── api/              # Backend server (Express, Go Chi, FastAPI, etc.)
├── packages/
│   └── db/               # Shared DB connection client (Prisma / Drizzle)
├── docker-compose.yml    # Local Database container setup
├── koko.config.json      # Koko workspace configuration manifest
├── package.json          # Root configuration with workspaces
├── turbo.json            # Turborepo pipeline configuration
└── README.md
```

---

## ⚙️ Configuration Manifest (`koko.config.json`)

On initialization, Koko outputs a manifest file at the root, capturing the active stack. This configuration drives future workspace operations:

```json
{
  "$schema": "https://koko-cli.dev/schema.json",
  "project": {
    "name": "my-super-app",
    "cliVersion": "v2.0.0",
    "createdAt": "2026-08-25T20:50:00Z"
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
    }
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

Contributions are welcome! Please check out our [Contributing Guide](CONTRIBUTING.md) to get started.

## 📄 License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
