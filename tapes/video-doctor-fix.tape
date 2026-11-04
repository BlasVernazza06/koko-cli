# Charm VHS Tape for Koko CLI - Video 3: koko doctor --fix
Output ../Koko/public/videos/koko-doctor-fix.mp4
Output ../Koko/public/videos/koko-doctor-fix.gif

Set Shell "powershell"
Set FontSize 14
Set FontFamily "JetBrains Mono, Fira Code, Cascadia Code, monospace"
Set Width 1120
Set Height 680
Set Padding 18
Set Theme "Catppuccin Mocha"
Set TypingSpeed 35ms

# Ensure base project exists
Hide
Type "if (!(Test-Path saas-app)) { .\koko.exe init saas-app -r saas }; Clear-Host"
Enter
Sleep 600ms
Show

# 1. Run doctor with --fix
Type ".\koko.exe doctor -d ./saas-app --fix"
Sleep 500ms
Enter

# 2. Wait for diagnostics & reconciliation output
Sleep 4000ms

# 3. Verify manifest
Type "Get-Content ./saas-app/koko.config.json -Head 15"
Sleep 400ms
Enter

Sleep 3000ms
