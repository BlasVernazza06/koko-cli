---
description: Scaffold a new Koko CLI fullstack project or monorepo. Plans a validated stack, confirms it with you, then generates it.
argument-hint: "<project name and/or stack preferences, e.g. my-app --recipe saas or nextjs express drizzle stripe>"
---

# Create a Koko CLI project

User request: $ARGUMENTS

Read the project name and stack preferences from the request above. If the request is empty, ask what kind of project the user wants to build.

Follow the **scaffold-project** skill:
1. **Settle the stack**: Parse framework, backend, database, ORM, auth, payments, email, and addons.
2. **Validate compatibility**: Verify rules (e.g. tRPC requires React, Mongoose requires MongoDB, etc.).
3. **Confirm with user**: Show the architecture plan.
4. **Generate**: Run `koko init <name> [flags]` with non-interactive flags.
5. **Verify & Report**: Inspect the generated files and report next steps (`pnpm install`, `pnpm dev`).
