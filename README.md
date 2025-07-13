# passgen

A secure, customizable password generator CLI tool written in Go.

## Features

- **Secure Random Generation**: Uses cryptographically secure random number generation
- **Customizable Character Sets**: Choose from lowercase, uppercase, numbers, and symbols
- **Advanced Options**: Exclude similar characters, custom character exclusion
- **Multiple Passwords**: Generate multiple passwords at once
- **Password Strength Checker**: Analyze password strength and get improvement suggestions
- **Preset Configurations**: Quick presets for common use cases
- **Cross-Platform**: Works on Linux, macOS, and Windows

## Installation

### Method 1: Automated Installation (Recommended)

The easiest way to install passgen with automatic PATH configuration:

```bash
# Download and run the installation script
curl -sSL https://raw.githubusercontent.com/kumarasakti/passgen/main/install.sh | bash

# Or clone and run locally
git clone https://github.com/kumarasakti/passgen.git
cd passgen
./install.sh
```

This script will:
- Install passgen using `go install`
- Automatically detect your shell (zsh, bash, fish)
- Add Go's bin directory to your PATH
- Make passgen immediately available

> **Note**: For PowerShell users, see [Method 3: PowerShell Installation](#method-3-powershell-installation-windowslinuxmacos) for a native PowerShell experience.

### Method 2: Quick One-liner

```bash
curl -sSL https://raw.githubusercontent.com/kumarasakti/passgen/main/quick-install.sh | bash
```

### Method 3: PowerShell Installation (Windows/Linux/macOS)

For PowerShell users on any platform:

```powershell
# Automated installation with interactive prompts
irm https://raw.githubusercontent.com/kumarasakti/passgen/main/install.ps1 | iex

# Quick installation without prompts
irm https://raw.githubusercontent.com/kumarasakti/passgen/main/quick-install.ps1 | iex

# Or run locally
git clone https://github.com/kumarasakti/passgen.git
cd passgen
.\install.ps1
```

This PowerShell script will:
- Install passgen using `go install`
- Automatically configure your PowerShell profile
- Add Go's bin directory to your PATH
- Work on Windows PowerShell, PowerShell Core, and Linux/macOS

**For Linux/macOS with PowerShell:**
```bash
curl -sSL https://raw.githubusercontent.com/kumarasakti/passgen/main/install.ps1 | pwsh
```

### Method 4: Manual Go Install

```bash
go install github.com/kumarasakti/passgen@latest
```

**Important**: After manual installation, you need to add Go's bin directory to your PATH:

```bash
# For zsh users
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
source ~/.zshrc

# For bash users
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc
```

**For PowerShell users:**
```powershell
# Add to your PowerShell profile
$env:PATH += ";$(go env GOPATH)\bin"  # Windows
$env:PATH += ":$(go env GOPATH)/bin"  # Linux/macOS
```

### Method 5: Download Pre-built Binaries

Download the latest release for your platform from the [releases page](https://github.com/kumarasakti/passgen/releases) and place it in your PATH.

### Method 6: Build from Source

```bash
git clone https://github.com/kumarasakti/passgen.git
cd passgen
make build
sudo mv build/passgen /usr/local/bin/  # Install to system PATH
```

### Verify Installation

```bash
passgen --version
passgen --help
```

If you get "command not found", your Go bin directory is not in PATH. Use Method 1 (automated installation) or follow the PATH setup instructions in Method 3.

## Usage

### Basic Usage

```bash
# Generate a default password (12 characters, letters and numbers)
passgen

# Generate a 16-character password with symbols
passgen -l 16 -s

# Generate 5 passwords at once
passgen -c 5

# Generate a secure password (all character types)
passgen --secure
```

### Advanced Options

```bash
# Exclude similar characters (i, l, 1, L, o, 0, O)
passgen --exclude-similar

# Exclude specific characters
passgen --exclude "aeiou"

# Generate only numbers (PIN)
passgen --lower=false --upper=false --numbers=true --symbols=false -l 6

# Generate alphanumeric password
passgen --alphanumeric -l 12
```

### Preset Configurations

```bash
# Secure password (16 chars, all types)
passgen preset secure

# Simple password (12 chars, letters and numbers)
passgen preset simple

# PIN (6 digits)
passgen preset pin

# Alphanumeric (12 chars, letters and numbers)
passgen preset alphanumeric
```

### Password Strength Checker

```bash
# Check password strength
passgen check "mypassword123"

# Example output:
# Strength: Medium (Score: 4/8)
# Suggestions: Add uppercase letters, Add special characters
```

## Command Line Options

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--length` | `-l` | Password length | 12 |
| `--count` | `-c` | Number of passwords to generate | 1 |
| `--lower` | | Include lowercase letters | true |
| `--upper` | | Include uppercase letters | true |
| `--numbers` | `-n` | Include numbers | true |
| `--symbols` | `-s` | Include symbols | false |
| `--exclude-similar` | | Exclude similar characters (il1Lo0O) | false |
| `--exclude` | | Characters to exclude | "" |
| `--secure` | `-S` | Generate secure password (all types) | false |
| `--simple` | `-m` | Generate simple password (letters + numbers) | false |
| `--alphanumeric` | `-a` | Generate alphanumeric password | false |
| `--help` | `-h` | Show help message | |
| `--version` | `-v` | Show version | |

## Examples

### Generate Different Types of Passwords

```bash
# Default password (14 characters, uppercase, lowercase, symbols)
$ passgen
🎉 Password Generated Successfully! 🎉

Password: -PUZ:@)x$DSloj

📊 Password Analysis:
✅ Length: 14 characters (Good!)
✅ Character Sets: Lowercase, Uppercase, Symbols
✅ Entropy: 88.0 bits (Very Strong!)
✅ Strength: Very Strong 💪

🔒 Security Assessment:
• This password would take approximately 4892016 years to crack with modern hardware
• Contains 3 different character types for excellent complexity
• Exceeds security standards for high-value accounts

Someone's taking this security thing seriously! 🌟

# Secure password with all character types (gets the ice cold message!)
$ passgen --secure -l 20
🎉 Password Generated Successfully! 🎉

Password: =;1.7y$A]RqH8a7):s&C

📊 Password Analysis:
✅ Length: 20 characters (Excellent!)
✅ Character Sets: Lowercase, Uppercase, Numbers, Symbols
✅ Entropy: 129.2 bits (Extremely Strong!)
✅ Strength: Extremely Strong 🔥

🔒 Security Assessment:
• This password would take approximately 1.2e+19 years to crack with modern hardware
• Contains 4 different character types for maximum complexity
• Quantum-resistant for the foreseeable future!

Brr, that's ice cold security! Even hackers are shivering! 🥶

# Medium strength password (with sarcastic feedback)
$ passgen -l 8
✨ Password Generated! ✨

Password: gz*yK#<&

📊 Password Analysis:
⚠️  Length: 8 characters (Could be longer)
✅ Character Sets: Lowercase, Uppercase, Symbols
✅ Entropy: 50.3 bits (Medium!)
✅ Strength: Medium ⚡

🔒 Security Assessment:
• This password would take approximately 11.4 minutes to crack with modern hardware
• Contains 3 different character types for excellent complexity
• Adequate for most general purposes

💡 Tips for improvement:
• Consider using 12+ characters for better security

Well, it's... adequate. I guess that's something! 👍
```

### Advanced Usage

```bash
# Generate multiple passwords
$ passgen -c 2 -l 10
🎉 Password Generated Successfully! 🎉

Password: CIMc|OVmzv

📊 Password Analysis:
✅ Length: 10 characters (Good!)
✅ Character Sets: Lowercase, Uppercase, Symbols
✅ Entropy: 62.9 bits (Strong!)
✅ Strength: Strong 💯

🔒 Security Assessment:
• This password would take approximately 1.2 years to crack with modern hardware
• Contains 3 different character types for excellent complexity
• Great for securing important accounts

Not bad, you actually read the security guidelines! 🎯

==================================================

🎉 Password Generated Successfully! 🎉

Password: oXEduV)%|K

📊 Password Analysis:
✅ Length: 10 characters (Good!)
✅ Character Sets: Lowercase, Uppercase, Symbols
✅ Entropy: 62.9 bits (Strong!)
✅ Strength: Strong 💯

🔒 Security Assessment:
• This password would take approximately 1.2 years to crack with modern hardware
• Contains 3 different character types for excellent complexity
• Great for securing important accounts

Not bad, you actually read the security guidelines! 🎯

# Check version
$ passgen --version
passgen version v1.0.3
```

## Development

### Installation Scripts

This repository includes multiple installation scripts for different environments:

- **`install.sh`** - Full-featured bash installer with shell detection (zsh, bash, fish)
- **`quick-install.sh`** - Minimal bash installer for quick setup
- **`install.ps1`** - Full-featured PowerShell installer (Windows/Linux/macOS)
- **`quick-install.ps1`** - Minimal PowerShell installer for quick setup

All scripts automatically handle PATH configuration and provide colored output with error handling.

### Building

```bash
# Install dependencies
make deps

# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Development workflow (format, lint, test, build)
make dev
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -run TestGeneratePassword
```

### Release

```bash
# Create release binaries and archives
make release
```

## Security

This tool uses Go's `crypto/rand` package for cryptographically secure random number generation. The generated passwords are suitable for production use.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Changelog

### v1.0.0
- Initial release
- Basic password generation with customizable options
- Password strength checker
- Preset configurations
- Cross-platform support
