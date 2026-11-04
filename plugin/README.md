# 🤖 Koko CLI AI Agent Plugin

Official AI Assistant Plugin for **Koko CLI** (compatible with Claude Code, Antigravity, Cursor, Codex, and other AI coding agents).

This plugin equips AI agents with the domain knowledge, rules, and commands needed to plan, scaffold, and extend modern fullstack TypeScript monorepos using Koko CLI rather than writing boilerplate by hand.

---

## ⚡ Installation

Install directly inside your AI assistant (e.g. Claude Code):

```bash
# Add the Koko CLI marketplace
/plugin marketplace add BlasVernazza06/koko-cli

# Install the Koko plugin
/plugin install koko@koko-cli
```

---

## 🚀 Capabilities & Commands

### 1. Natural Language Scaffolding
Prompt your agent in natural language:
> *"Create a fullstack SaaS with Next.js, Express, Postgres, Drizzle, Better-Auth, and Stripe"*

The bundled **Stack Architect** agent evaluates stack rules, presents an architectural blueprint, and scaffolds the monorepo via `koko init`.

### 2. Slash Commands

| Command | Usage | Description |
| :--- | :--- | :--- |
| **`/koko:new`** | `/koko:new <name> [options]` | Plans and scaffolds a new fullstack project or recipe. |
| **`/koko:add`** | `/koko:add <addon>` | Injects services (Stripe, Polar, Resend, Brevo, shadcn, Lucide, Docker) into an existing app. |
| **`/koko:doctor`** | `/koko:doctor [--fix]` | Analyzes project health and reconciles `koko.config.json` discrepancies. |

---

## 📂 Plugin Structure

```text
plugin/
├── .claude-plugin/
│   └── plugin.json          # Plugin manifest and metadata
├── README.md                # Plugin documentation
├── commands/                # Interactive slash commands
│   ├── new.md               # /koko:new
│   ├── add.md               # /koko:add
│   └── doctor.md            # /koko:doctor
├── agents/                  # Specialist subagents
│   └── stack-architect.md   # Architectural design and validation agent
└── skills/                  # Domain expertise for AI assistants
    ├── scaffold-project/    # Project creation workflows and compatibility rules
    ├── add-to-project/      # Service injection workflows
    └── doctor-project/      # Diagnostic and auto-repair workflows
```

---

## 📖 Learn More

- **Repository**: [https://github.com/BlasVernazza06/koko-cli](https://github.com/BlasVernazza06/koko-cli)
- **License**: MIT
