# Script to cross-compile Go binaries for all platforms into the NPM folder

$NpmBinDir = Join-Path (Get-Location) "npm/bin"

# Create destination folder if it doesn't exist
if (-not (Test-Path $NpmBinDir)) {
    New-Item -ItemType Directory -Force -Path $NpmBinDir | Out-Null
}

Write-Host "Compiling Koko CLI binaries..." -ForegroundColor Cyan

# 1. Windows amd64
Write-Host "-> Building for Windows (amd64)..." -ForegroundColor Yellow
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -ldflags="-s -w" -o (Join-Path $NpmBinDir "koko-windows-amd64.exe") main.go

# 2. macOS amd64
Write-Host "-> Building for macOS (amd64)..." -ForegroundColor Yellow
$env:GOOS = "darwin"
$env:GOARCH = "amd64"
go build -ldflags="-s -w" -o (Join-Path $NpmBinDir "koko-darwin-amd64") main.go

# 3. macOS arm64 (M1/M2/M3)
Write-Host "-> Building for macOS (arm64)..." -ForegroundColor Yellow
$env:GOOS = "darwin"
$env:GOARCH = "arm64"
go build -ldflags="-s -w" -o (Join-Path $NpmBinDir "koko-darwin-arm64") main.go

# 4. Linux amd64
Write-Host "-> Building for Linux (amd64)..." -ForegroundColor Yellow
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -ldflags="-s -w" -o (Join-Path $NpmBinDir "koko-linux-amd64") main.go

# Reset environment variables
$env:GOOS = ""
$env:GOARCH = ""

Write-Host "Build complete! Binaries are located in: $NpmBinDir" -ForegroundColor Green
