#!/usr/bin/env pwsh

# Quick passgen installer - PowerShell one-liner version
# Usage: 
#   irm https://raw.githubusercontent.com/kumarasakti/passgen/main/quick-install.ps1 | iex
#   curl -sSL https://raw.githubusercontent.com/kumarasakti/passgen/main/quick-install.ps1 | pwsh

# Color functions
function Write-Success { param($Message) Write-Host "✅ $Message" -ForegroundColor Green }
function Write-Info { param($Message) Write-Host "🚀 $Message" -ForegroundColor Cyan }
function Write-Error { param($Message) Write-Host "❌ $Message" -ForegroundColor Red }
function Write-Celebrate { param($Message) Write-Host "🎉 $Message" -ForegroundColor Magenta }

Write-Info "Quick installing passgen..."

# Check if Go is installed
try {
    $null = go version 2>$null
    if ($LASTEXITCODE -ne 0) { throw }
} catch {
    Write-Error "Go is not installed. Please install Go from https://golang.org/dl/"
    exit 1
}

# Install passgen
try {
    go install github.com/kumarasakti/passgen@latest
    if ($LASTEXITCODE -ne 0) { throw }
} catch {
    Write-Error "Failed to install passgen"
    exit 1
}

# Get Go bin path
$goPath = go env GOPATH
$goBin = go env GOBIN
if (-not $goBin -or $goBin -eq "") {
    $goBin = if ($IsWindows -or $env:OS -eq "Windows_NT") {
        Join-Path $goPath "bin"
    } else {
        "$goPath/bin"
    }
}

# Add to current session PATH
if ($IsWindows -or $env:OS -eq "Windows_NT") {
    $env:PATH += ";$goBin"
} else {
    $env:PATH += ":$goBin"
}

# Add to PowerShell profile
$profilePath = $PROFILE.CurrentUserCurrentHost
$pathExport = "`$env:PATH += `"$(if ($IsWindows -or $env:OS -eq 'Windows_NT') {';'} else {':'})`$(go env GOPATH)$(if ($IsWindows -or $env:OS -eq 'Windows_NT') {'\\bin'} else {'/bin'})`""

# Ensure profile directory exists
$profileDir = Split-Path $profilePath -Parent
if (-not (Test-Path $profileDir)) {
    New-Item -ItemType Directory -Path $profileDir -Force | Out-Null
}

# Check if already in profile
$profileContent = if (Test-Path $profilePath) { Get-Content $profilePath -Raw } else { "" }
if ($profileContent -notmatch "go env GOPATH") {
    Add-Content -Path $profilePath -Value $pathExport
    Write-Success "Added Go bin to PATH in PowerShell profile"
}

Write-Celebrate "passgen installed! Try: passgen --version"
