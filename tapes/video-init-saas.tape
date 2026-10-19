# Charm VHS Tape for Koko CLI - Video 1: Quick Blueprint (saas-app)
Output ../Koko/public/videos/koko-init-saas.mp4
Output ../Koko/public/videos/koko-init-saas.gif

Set Shell "powershell"
Set FontSize 14
Set FontFamily "JetBrains Mono, Fira Code, Cascadia Code, monospace"
Set Width 1120
Set Height 680
Set Padding 18
Set Theme "Catppuccin Mocha"
Set TypingSpeed 35ms

# Prepare clean environment
Hide
Type "Remove-Item -Recurse -Force saas-app -ErrorAction SilentlyContinue; Clear-Host"
Enter
Sleep 600ms
Show

# 1. Direct Command execution with Recipe
Type ".\koko.exe init saas-app -r saas"
Sleep 500ms
Enter

# 2. Wait for scaffolding execution & celebration
Sleep 4000ms

# 3. Inspect generated structure
Type "Get-ChildItem saas-app | Select-Object Name"
Sleep 400ms
Enter

Sleep 3000ms
