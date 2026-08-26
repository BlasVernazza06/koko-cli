# Tape file for Charm VHS - generates the terminal walkthrough GIF
# Output path
Output demo.gif

# Set terminal configurations
Set Shell "powershell"
Set FontSize 16
Set Width 1200
Set Height 600
Set Padding 20
Set Theme "Catppuccin Mocha"

# Start typing the command
Type "koko init"
Sleep 500ms
Enter

# Wait for project name prompt
Sleep 1s
Type "koko-monorepo"
Sleep 500ms
Enter

# Choose Quick Setup vs Manual Configuration
Sleep 1s
Down
Sleep 800ms
Enter

# Manual options selection (Step-by-step)
# 1. Frontend: Next.js
Sleep 1s
Enter

# 2. Backend: Go / Chi Router
Sleep 1s
Down
Down
Sleep 500ms
Enter

# 3. Package Manager: Go Modules (auto-selected / pnpm)
Sleep 1s
Enter

# 4. Database: PostgreSQL
Sleep 1s
Enter

# 5. ORM: Drizzle ORM (Prisma/drizzle compatible)
Sleep 1s
Enter

# 6. Addons: Docker Compose
Sleep 1s
Enter

# 7. Git: Yes
Sleep 1s
Enter

# Wait for scaffolding completion
Sleep 4s

# End recording
Sleep 2s
