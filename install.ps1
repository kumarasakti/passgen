#!/usr/bin/env pwsh

# passgen Installation Script for PowerShell
# This script installs passgen and automatically configures your PATH
# Compatible with Windows PowerShell, PowerShell Core, and Linux/macOS

param(
    [switch]$QuickInstall,
    [switch]$Help
)

# Color functions for better output
function Write-Success { param($Message) Write-Host "✅ $Message" -ForegroundColor Green }
function Write-Info { param($Message) Write-Host "🚀 $Message" -ForegroundColor Cyan }
function Write-Warning { param($Message) Write-Host "⚠️  $Message" -ForegroundColor Yellow }
function Write-Error { param($Message) Write-Host "❌ $Message" -ForegroundColor Red }
function Write-Celebrate { param($Message) Write-Host "🎉 $Message" -ForegroundColor Magenta }

function Show-Help {
    Write-Host @"
passgen PowerShell Installer

USAGE:
    .\install.ps1 [OPTIONS]

OPTIONS:
    -QuickInstall    Perform a quick installation without interactive prompts
    -Help           Show this help message

EXAMPLES:
    .\install.ps1                    # Interactive installation
    .\install.ps1 -QuickInstall      # Quick installation
    
REMOTE INSTALLATION:
    # Windows PowerShell / PowerShell Core
    irm https://raw.githubusercontent.com/kumarasakti/passgen/main/install.ps1 | iex
    
    # Linux/macOS with PowerShell
    curl -sSL https://raw.githubusercontent.com/kumarasakti/passgen/main/install.ps1 | pwsh

"@ -ForegroundColor White
    exit 0
}

if ($Help) {
    Show-Help
}

Write-Info "Installing passgen..."

# Check if Go is installed
try {
    $goVersion = go version 2>$null
    if (-not $goVersion) {
        throw "Go not found"
    }
    Write-Success "Go is installed: $($goVersion.Split(' ')[2])"
} catch {
    Write-Error "Go is not installed or not in PATH."
    Write-Host "Please install Go from https://golang.org/dl/" -ForegroundColor Yellow
    exit 1
}

# Install passgen using go install
Write-Info "Installing passgen from GitHub..."
try {
    go install github.com/kumarasakti/passgen@latest
    if ($LASTEXITCODE -ne 0) {
        throw "go install failed"
    }
} catch {
    Write-Error "Failed to install passgen. Please check your internet connection and Go installation."
    exit 1
}

# Get Go paths
try {
    $goPath = go env GOPATH
    $goBin = go env GOBIN
    
    # If GOBIN is not set, use GOPATH/bin
    if (-not $goBin -or $goBin -eq "") {
        if ($IsWindows -or $env:OS -eq "Windows_NT") {
            $goBin = Join-Path $goPath "bin"
        } else {
            $goBin = "$goPath/bin"
        }
    }
    
    Write-Success "Go bin directory: $goBin"
} catch {
    Write-Error "Failed to get Go environment variables."
    exit 1
}

# Check if passgen binary exists
$passggenBinary = if ($IsWindows -or $env:OS -eq "Windows_NT") {
    Join-Path $goBin "passgen.exe"
} else {
    "$goBin/passgen"
}

if (-not (Test-Path $passggenBinary)) {
    Write-Error "Installation failed. Binary not found at: $passggenBinary"
    exit 1
}

Write-Success "passgen binary installed to: $passggenBinary"

# Function to check if directory is in PATH
function Test-InPath {
    param([string]$Directory)
    
    $pathSeparator = if ($IsWindows -or $env:OS -eq "Windows_NT") { ";" } else { ":" }
    $currentPath = $env:PATH -split $pathSeparator
    
    return $currentPath -contains $Directory
}

# Function to get PowerShell profile path
function Get-ProfilePath {
    # Check for different PowerShell profiles in order of preference
    $profiles = @(
        $PROFILE.CurrentUserAllHosts,
        $PROFILE.CurrentUserCurrentHost,
        $PROFILE.AllUsersCurrentHost
    )
    
    foreach ($profile in $profiles) {
        if ($profile -and (Test-Path (Split-Path $profile -Parent) -ErrorAction SilentlyContinue)) {
            return $profile
        }
    }
    
    # Default to CurrentUserCurrentHost
    return $PROFILE.CurrentUserCurrentHost
}

# Function to configure PATH
function Set-PassgenPath {
    $pathConfigured = $false
    
    if (Test-InPath $goBin) {
        Write-Success "Go bin directory is already in your PATH"
        return $true
    }
    
    Write-Info "Configuring PATH..."
    
    # Get PowerShell profile path
    $profilePath = Get-ProfilePath
    Write-Host "Profile path: $profilePath" -ForegroundColor Gray
    
    # Ensure profile directory exists
    $profileDir = Split-Path $profilePath -Parent
    if (-not (Test-Path $profileDir)) {
        try {
            New-Item -ItemType Directory -Path $profileDir -Force | Out-Null
            Write-Success "Created profile directory: $profileDir"
        } catch {
            Write-Warning "Could not create profile directory: $profileDir"
            return $false
        }
    }
    
    # Check if PATH configuration already exists
    $pathExport = "`$env:PATH += `"$(if ($IsWindows -or $env:OS -eq 'Windows_NT') {';'} else {':'})`$(go env GOPATH)$(if ($IsWindows -or $env:OS -eq 'Windows_NT') {'\\bin'} else {'/bin'})`""
    
    if ((Test-Path $profilePath) -and (Get-Content $profilePath -Raw -ErrorAction SilentlyContinue) -match "go env GOPATH") {
        Write-Warning "Go PATH configuration already exists in profile"
    } else {
        try {
            # Add PATH configuration to profile
            Add-Content -Path $profilePath -Value "`n# Added by passgen installer"
            Add-Content -Path $profilePath -Value $pathExport
            Write-Success "Added Go bin to PATH in PowerShell profile"
            $pathConfigured = $true
        } catch {
            Write-Warning "Could not modify PowerShell profile: $profilePath"
        }
    }
    
    # Add to current session PATH
    if ($IsWindows -or $env:OS -eq "Windows_NT") {
        $env:PATH += ";$goBin"
    } else {
        $env:PATH += ":$goBin"
    }
    
    return $true
}

# Configure PATH
$pathResult = Set-PassgenPath

# Test if passgen is accessible
try {
    $version = & passgen --version 2>$null
    if ($LASTEXITCODE -eq 0) {
        Write-Host ""
        Write-Celebrate "Installation successful!"
        Write-Success "passgen is now available in your PATH"
        Write-Host "Version: $version" -ForegroundColor Gray
        Write-Host ""
        Write-Host "Try it out:" -ForegroundColor Yellow
        Write-Host "  passgen --version" -ForegroundColor White
        Write-Host "  passgen --help" -ForegroundColor White
        Write-Host "  passgen" -ForegroundColor White
        Write-Host ""
    } else {
        throw "passgen command failed"
    }
} catch {
    Write-Host ""
    Write-Warning "Installation completed but passgen is not immediately available."
    Write-Host "Please restart your PowerShell session or run:" -ForegroundColor Yellow
    Write-Host "  . `$PROFILE" -ForegroundColor White
    Write-Host ""
    Write-Host "Or run passgen with full path:" -ForegroundColor Yellow
    Write-Host "  & '$passggenBinary'" -ForegroundColor White
    Write-Host ""
}

Write-Host "📖 For more information, visit: " -NoNewline -ForegroundColor Gray
Write-Host "https://github.com/kumarasakti/passgen" -ForegroundColor Blue

# If this script was downloaded and executed directly, offer to clean up
if ($MyInvocation.MyCommand.Path -and $MyInvocation.MyCommand.Path -match "tmp|temp") {
    Write-Host ""
    if (-not $QuickInstall) {
        $cleanup = Read-Host "Remove installer script? (y/N)"
        if ($cleanup -eq "y" -or $cleanup -eq "Y") {
            try {
                Remove-Item $MyInvocation.MyCommand.Path -Force
                Write-Success "Installer script removed"
            } catch {
                Write-Warning "Could not remove installer script"
            }
        }
    }
}
