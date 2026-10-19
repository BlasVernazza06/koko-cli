---
description: Run stack diagnostics on the current Koko project to detect discrepancies against koko.config.json, optionally fixing them automatically.
argument-hint: "[--fix / -f]"
---

# Diagnose and Repair Koko Project

User request: $ARGUMENTS

Follow the **doctor-project** skill:
1. Check if `koko.config.json` exists in the current directory.
2. Run `koko doctor` to analyze project structure, dependencies, ORMs, and services.
3. If `--fix` or `-f` is requested, run `koko doctor --fix` to reconcile differences and synchronize `koko.config.json`.
4. Report detected differences and resolution status.
