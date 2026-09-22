# Charm VHS Tape for Koko CLI Demo (Manual Customization Mode)
Output assets/demo.gif
Output assets/demo.mp4

Set Shell "powershell"
Set FontSize 14
Set FontFamily "JetBrains Mono, Fira Code, Cascadia Code, monospace"
Set Width 1120
Set Height 680
Set Padding 18
Set Theme "Catppuccin Mocha"
Set TypingSpeed 35ms

# Hide shell startup noise
Hide
Type "Clear-Host"
Enter
Sleep 500ms
Show

# 1. Run Koko CLI
Type ".\koko.exe init"
Sleep 400ms
Enter

# 2. Enter Project Name
Sleep 700ms
Type "my-custom-stack"
Sleep 500ms
Enter

# 3. Select Setup Mode: Navigate to Manual Configuration
Sleep 700ms
Down
Sleep 600ms
Enter

# 4. Step-by-Step Manual Selection
# 4.1 Frontend: Next.js (App Router)
Sleep 600ms
Down
Sleep 300ms
Up
Sleep 300ms
Enter

# 4.2 Backend: Hono / Node.js
Sleep 600ms
Down
Sleep 400ms
Enter

# 4.3 API Layer: tRPC
Sleep 600ms
Enter

# 4.4 Package Manager: PNPM
Sleep 600ms
Enter

# 4.5 Database: PostgreSQL
Sleep 600ms
Enter

# 4.6 ORM: Drizzle ORM
Sleep 600ms
Enter

# 4.7 Auth: Better-Auth
Sleep 600ms
Enter

# 4.8 Addons (Multi-select with Spacebar):
Sleep 600ms
Space
Sleep 300ms
Down
Sleep 200ms
Space
Sleep 300ms
Down
Sleep 200ms
Down
Sleep 200ms
Space
Sleep 300ms
Down
Sleep 200ms
Down
Sleep 200ms
Space
Sleep 400ms
Enter

# 4.9 Git: Initialize Git Repository
Sleep 600ms
Enter

# 5. Scaffolding in progress with live checklist
Sleep 3s

# 6. Finished screen with full stack summary & next steps
Sleep 4s


