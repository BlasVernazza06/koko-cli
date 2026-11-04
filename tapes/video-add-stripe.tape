# Charm VHS Tape for Koko CLI - Video 2: koko add stripe
Output ../Koko/public/videos/koko-add-stripe.mp4
Output ../Koko/public/videos/koko-add-stripe.gif

Set Shell "powershell"
Set FontSize 14
Set FontFamily "JetBrains Mono, Fira Code, Cascadia Code, monospace"
Set Width 1120
Set Height 680
Set Padding 18
Set Theme "Catppuccin Mocha"
Set TypingSpeed 35ms

# Prepare clean environment: scaffold a clean monorepo without payments, then enter the project
Hide
Type "Remove-Item -Recurse -Force saas-app -ErrorAction SilentlyContinue; $env:PATH = '$PWD;' + $env:PATH; Set-Alias koko (Join-Path $PWD 'koko.exe'); koko init saas-app -f next -b express --database postgres --orm drizzle --auth better-auth -p pnpm --git no > $null; Set-Location saas-app; Clear-Host"
Enter
Sleep 1200ms
Show

# 1. Run koko add stripe directly inside the project (as in the guide)
Type "koko add stripe"
Sleep 600ms
Enter

# 2. Wait for injection header, VFS processing and confirmation
Sleep 3200ms

# 3. View injected environment variables in .env
Type "Get-Content .env | Select-String 'STRIPE'"
Sleep 500ms
Enter

# 4. Show that package.json was updated with stripe
Sleep 1800ms
Type "Get-Content apps/api/package.json | Select-String 'stripe'"
Sleep 500ms
Enter

Sleep 3000ms
