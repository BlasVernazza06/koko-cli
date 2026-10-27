---
name: turborepo
description: High-performance monorepo orchestration, task pipelines, caching dependencies, packages boundaries, remote caching, and workspace protocols.
globs: ["turbo.json", "package.json", "pnpm-workspace.yaml"]
---

# Turborepo Monorepo Architecture & Pipeline Standards

Turborepo orchestrates task execution, dependency graphs, and multi-package caching across modern TypeScript/JavaScript monorepos.

## Core Concepts & Task Dependency Graph

- **`dependsOn`**: Declare execution order. `^build` guarantees all workspace dependencies build before the current package builds.
- **`outputs`**: List of filesystem artifacts saved in local (`.turbo/cache`) and remote cache.
- **`inputs`**: Files that trigger cache invalidation when modified. Omit lockfiles and docs from task inputs to prevent unnecessary cache misses.
- **`cache: false`**: Required for dev servers and watch tasks (`dev`, `start`) to prevent daemon hanging.

## Quick Reference Table: Pipeline Directives

| Pipeline Key | Type | Purpose | Production Example |
| :--- | :--- | :--- | :--- |
| `dependsOn` | `string[]` | Upstream task prerequisites | `["^build", "^codegen"]` |
| `inputs` | `string[]` | Files hash inputs for cache key | `["src/**", "package.json", "tsconfig.json"]` |
| `outputs` | `string[]` | Output folders to store/restore | `[".next/**", "!.next/cache/**", "dist/**"]` |
| `cache` | `boolean` | Enable or disable task caching | `true` (default), `false` (dev) |
| `persistent` | `boolean` | Long-running interactive processes | `true` for watch/dev processes |
| `env` | `string[]` | Env vars included in cache hash | `["NODE_ENV", "NEXT_PUBLIC_*"]` |
| `globalEnv` | `string[]` | Global environment dependencies | `["DATABASE_URL", "VERCEL_ENV"]` |

## Pipeline Configuration Example (`turbo.json`)

```json
{
  "$schema": "https://turbo.build/schema.json",
  "ui": "tui",
  "globalDependencies": ["**/.env.*local"],
  "tasks": {
    "build": {
      "dependsOn": ["^build"],
      "inputs": ["src/**", "package.json", "tsconfig.json"],
      "outputs": [".next/**", "!.next/cache/**", "dist/**"]
    },
    "lint": {
      "dependsOn": ["^lint"]
    },
    "test": {
      "dependsOn": ["^build"],
      "inputs": ["src/**", "test/**"]
    },
    "dev": {
      "cache": false,
      "persistent": true
    }
  }
}
```

## Anti-Pattern vs Recommended Solution

### ❌ Anti-Pattern: Relative Cross-Package Imports
Directly importing source files across package folders breaks encapsulation, skips the build pipeline, and ruins cache isolation.

```typescript
// apps/web/src/app/page.tsx - BAD: Direct relative import escaping package boundary
import { UserAvatar } from '../../packages/ui/src/UserAvatar';
import { db } from '../../packages/db/src/client';
```

### ✅ Recommended: Workspace Package Boundary Imports
Expose typed exports in `packages/ui/package.json` and consume via workspace protocol.

```json
// apps/web/package.json
{
  "dependencies": {
    "@repo/ui": "workspace:*",
    "@repo/db": "workspace:*"
  }
}
```

```typescript
// apps/web/src/app/page.tsx - GOOD: Clean workspace alias
import { UserAvatar } from '@repo/ui/components/user-avatar';
import { db } from '@repo/db';
```

## Monorepo Architecture Layout

```
monorepo/
├── apps/
│   ├── web/                     # Next.js / Astro frontend application
│   └── api/                     # Node.js / Hono backend service
├── packages/
│   ├── tsconfig/                # Centralized tsconfig base configurations
│   ├── ui/                      # Shared design system components & tokens
│   ├── db/                      # Database clients, schemas, migrations
│   └── auth/                    # Shared authentication contracts
├── turbo.json                   # Pipeline orchestrator
└── pnpm-workspace.yaml          # Workspaces definition
```

## Remote Caching & CI Filtering

```bash
# Authenticate and link remote cache (Vercel or custom self-hosted cache)
npx turbo login
npx turbo link

# Run task only on changed packages and their dependents
turbo run build --filter=...[origin/main]

# Run tests only on affected apps
turbo run test --filter=web...
```
