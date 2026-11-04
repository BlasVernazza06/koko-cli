---
name: add-to-project
description: Add services, dependencies or tooling addons to an existing Koko CLI project, such as Stripe, Polar, Resend, Brevo, shadcn/ui, Lucide, SVGL, Framer Motion, Zod, Docker Compose, or GitHub Actions CI. Use when the project already has a koko.config.json and the user wants to extend or add features.
metadata:
  priority: 8
  docs:
    - "https://github.com/BlasVernazza06/koko-cli"
  pathPatterns:
    - "koko.config.json"
---

# Add to an existing Koko CLI project

Inject new services and dependencies into an existing project with `koko add` instead of wiring `package.json` dependencies and `.env` credentials by hand.

## Supported Addons & Services

| Addon Argument | Category | What Koko Adds |
| :--- | :--- | :--- |
| **`stripe`** | Payments | Injects `stripe` package in `apps/api` and `apps/web`, adds `STRIPE_SECRET_KEY` and `NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY` in `.env`. |
| **`polar`** | Payments | Injects `@polar-sh/sdk`, adds `POLAR_ACCESS_TOKEN` and `POLAR_ORGANIZATION_ID` in `.env`. |
| **`resend`** | Email | Injects `resend`, adds `RESEND_API_KEY` in `.env`. |
| **`brevo`** | Email | Injects `@getbrevo/brevo`, adds `BREVO_API_KEY` in `.env`. |
| **`shadcn`** | UI Components | Injects `@<app>/ui` workspace dependency and copies `packages/ui`. |
| **`lucide`** | Icons | Injects `lucide-react` in `apps/web`. |
| **`svgl`** | Brand Logos | Injects `@ridemountainpig/svgl-react` in `apps/web`. |
| **`motion`** | Animations | Injects `framer-motion` in `apps/web`. |
| **`zod`** | Validation | Injects `zod` in `apps/api` and `apps/web`. |
| **`docker`** | DevOps | Copies and configures `docker-compose.yml` based on active database. |
| **`github_actions`** | CI/CD | Generates `.github/workflows/ci.yml`. |

## Steps

1. **Verify `koko.config.json`**:
   - Ensure the current working directory contains `koko.config.json`. If missing, notify the user that this command only runs inside an existing Koko project.

2. **Check Current Stack**:
   - Inspect `koko.config.json` to verify if the addon is already present (e.g. `features.payments` or `features.email`).

3. **Execute Ingestion**:
   ```bash
   koko add <addon>
   ```

4. **Verify Output**:
   - Confirm that `koko.config.json` was updated.
   - Inspect `.env` and `.env.example` to ensure new keys were appended cleanly.
   - Inspect `package.json` to verify workspace dependencies were resolved.
   - Inform the user of next steps (e.g. running `pnpm install` or providing API keys in `.env`).
