---
description: Add a new service, dependency or tooling addon to the current Koko project. Plans the change, confirms it, then injects it.
argument-hint: "<addon name: stripe, polar, resend, brevo, shadcn, lucide, svgl, motion, zod, docker, github_actions>"
---

# Add to an existing Koko CLI project

User request: $ARGUMENTS

Read the addon or service name from the request above. If the request is empty, ask what the user wants to add.

Follow the **add-to-project** skill:
1. Verify that `koko.config.json` exists in the project root.
2. Check if the addon is already installed.
3. Show the planned changes (manifest update, package dependencies, environment variables).
4. Run `koko add <addon>`.
5. Report changes made to `package.json` and `.env`.
