---
name: github-actions
description: Automated GitHub Actions CI/CD workflows, dependency caching, matrix builds, permissions hardening, OIDC security, and fast feedback loops.
globs: [".github/workflows/*.yml", ".github/workflows/*.yaml"]
---

# Production GitHub Actions CI/CD Standards

Automating build, test, and release pipelines with aggressive dependency caching, least-privilege permissions, and deterministic reproducibility.

## Workflow Permissions & Security Reference

| Permission Scope | Default Security Principle | Recommended Setting |
| :--- | :--- | :--- |
| `permissions: contents` | Principle of least privilege | `read` (read-only for checkout) |
| `permissions: packages` | Registry write access | `none` or `write` only on release jobs |
| `permissions: pull-requests`| PR comment integration | `write` only when posting test coverage |
| `permissions: id-token` | OIDC Cloud Provider auth (AWS/GCP/Cloudflare) | `write` only on cloud deployment steps |

## Performance & Optimization Rules

1. **Concurrency Control**: Automatically cancel previous in-progress runs when a new commit is pushed to the same PR.
2. **Action Pinning**: Pin all third-party actions to full commit SHAs or verified major tags (`actions/checkout@v4`).
3. **Cache Setup**: Use built-in cache mechanisms (`setup-node` with `cache: 'pnpm'`, `setup-go` with `cache: true`).
4. **Fast-Failing Pipeline Stages**: Organize workflows into parallel stages so quick failures (linting, type-checks) abort heavy integration jobs immediately.

## Anti-Pattern vs Recommended Solution

### ❌ Anti-Pattern: Unbound Permissions with Redundant Dependency Downloads
Running workflows without concurrency groups or caching results in 10-minute wasted builds and exposed token permissions.

```yaml
# BAD: Default broad permissions, missing caching, duplicate concurrent runs
name: CI
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: 22
      - run: npm install  # Slow: Downloads all dependencies from npm registry every single run!
      - run: npm test
```

### ✅ Recommended: Hardened Matrix Workflow with Concurrency & Caching

```yaml
name: CI & Quality Gate

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

# Cancel redundant pending runs on the same branch or PR
concurrency:
  group: ${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true

# Hardening: Default all jobs to read-only access
permissions:
  contents: read

jobs:
  lint-and-types:
    name: Lint & Typecheck
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Install pnpm
        uses: pnpm/action-setup@v4
        with:
          version: 9

      - name: Setup Node.js with pnpm cache
        uses: actions/setup-node@v4
        with:
          node-version: 22
          cache: 'pnpm'

      - name: Install dependencies
        run: pnpm install --frozen-lockfile

      - name: Run lint
        run: pnpm run lint

      - name: Run typecheck
        run: pnpm run check-types

  test:
    name: Test Suite (${{ matrix.os }})
    needs: [lint-and-types]
    runs-on: ${{ matrix.os }}
    strategy:
      fail-fast: false
      matrix:
        os: [ubuntu-latest]
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Install pnpm
        uses: pnpm/action-setup@v4
        with:
          version: 9

      - name: Setup Node.js with cache
        uses: actions/setup-node@v4
        with:
          node-version: 22
          cache: 'pnpm'

      - name: Install dependencies
        run: pnpm install --frozen-lockfile

      - name: Run tests with coverage
        run: pnpm run test:coverage
```
