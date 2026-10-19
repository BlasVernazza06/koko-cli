---
name: doctor-project
description: Inspect and diagnose an existing project stack against koko.config.json, detect discrepancies, and auto-repair manifest state using koko doctor --fix. Use when auditing project health or synchronizing configuration after manual package changes.
metadata:
  priority: 7
  docs:
    - "https://github.com/BlasVernazza06/koko-cli"
---

# Diagnose and Repair Koko Projects

Use `koko doctor` to analyze the project's real state on disk (dependencies, frameworks, databases, and ORMs) and compare it against `koko.config.json`.

## When to Use

- When the user asks to "check the project", "verify the stack", or "fix koko config".
- When dependencies or database providers were changed manually in `package.json`.
- When `koko.config.json` is missing or out of sync.

## Steps

1. **Run Diagnostics**:
   ```bash
   koko doctor
   ```
   Inspect the output. Koko will categorize detected differences:
   - `+ [ADDED]`: Services or packages found on disk that aren't recorded in the config.
   - `~ [MODIFIED]`: Frameworks, ORMs or layouts that changed.
   - `- [REMOVED]`: Services previously recorded that are no longer present.

2. **Auto-Repair Discrepancies**:
   If differences exist or the user asks to fix the project state:
   ```bash
   koko doctor --fix
   ```
   Koko will write an updated, validated `koko.config.json` matching the real code.

3. **Report**:
   Summarize the changes applied and confirm that the project is now 100% in sync.
