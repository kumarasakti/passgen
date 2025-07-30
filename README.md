# passgen

A secure, customizable password generator CLI tool written in Go with clean architecture.

## 🎯 New in v1.1.0: Word-Based Password Generation!

Transform memorable words into secure passwords using intelligent transformation strategies.

```bash
# Transform a word into a secure password
passgen word "security"
# Output: S3cur1ty!42

# Choose transformation strategy
passgen word "password" --strategy leetspeak
# Output: p@ssw0rd

# Set complexity level
passgen word "secret" --complexity high --count 3
# Generate multiple variations
```

## Features

- **🎯 Word-Based Passwords**: Transform memorable words into secure passwords
- **🔧 Multiple Transformation Strategies**: Leetspeak, mixed-case, hybrid, prefix/suffix
- **📊 Complexity Levels**: Low, medium, high complexity transformations
- **🔐 Secure Random Generation**: Uses cryptographically secure random number generation
- **🎨 Customizable Character Sets**: Choose from lowercase, uppercase, numbers, and symbols
- **⚙️ Advanced Options**: Exclude similar characters, custom character exclusion
- **📦 Multiple Passwords**: Generate multiple passwords at once
- **🔍 Password Strength Checker**: Analyze password strength and get improvement suggestions
- **🚀 Preset Configurations**: Quick presets for common use cases
- **🌍 Cross-Platform**: Works on Linux, macOS, and Windows
- **🏗️ Clean Architecture**: Maintainable, testable, extensible codebase

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

### Quick Start

```bash
# Generate a default password (12 characters, letters and numbers)
passgen

# Generate a 16-character password with symbols
passgen -l 16 -s

# Transform a word into a secure password (NEW!)
passgen word "security"
# Output: S3cur1ty!42

# Generate 5 passwords at once
passgen -c 5

# Generate a secure password (all character types)
passgen --secure
```

### ⭐ Word-Based Password Generation (NEW in v1.1.0)

Transform memorable words into secure passwords using intelligent transformation strategies:

```bash
# Basic word transformation
passgen word "sunshine"
# Output: Sunsh1n3!42

# Use specific transformation strategies
passgen word "password" --strategy leetspeak
# Output: p@55w0rd

passgen word "myhouse" --strategy mixedcase
# Output: MyHoUsE

passgen word "admin" --strategy suffix
# Output: admin_789!

# Set complexity levels
passgen word "secret" --complexity low
# Output: Secret

passgen word "secret" --complexity medium
# Output: Secret42

passgen word "secret" --complexity high
# Output: S3cr3t!78$

# Generate multiple variations
passgen word "team" --count 5
# Output: Multiple secure variations of "team"
```

**Transformation Strategies:**
- `leetspeak`: Replace letters with numbers/symbols (a→@, e→3, i→1, o→0, s→$)
- `mixedcase`: Alternate between upper and lower case
- `suffix`: Add random numbers and symbols at the end
- `prefix`: Add random numbers and symbols at the beginning
- `insert`: Insert random characters throughout the word
- `hybrid`: Combine multiple strategies for maximum security (default)

**Complexity Levels:**
- `low`: Basic transformation (capitalization only)
- `medium`: Moderate transformation (some substitutions + numbers)
- `high`: Complex transformation (full substitutions + symbols + numbers)

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

### Standard Password Generation

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

### Word-Based Password Generation (NEW!)

```bash
passgen word <word> [flags]
```

| Flag | Short | Description | Default | Options |
|------|-------|-------------|---------|---------|
| `--strategy` | `-s` | Transformation strategy | hybrid | leetspeak, mixedcase, suffix, prefix, insert, hybrid |
| `--complexity` | `-x` | Complexity level | medium | low, medium, high |
| `--count` | `-c` | Number of variations to generate | 1 | |

**Strategy Details:**
- **leetspeak**: `security` → `s3cur1ty` (replaces a→@, e→3, i→1, o→0, s→$)
- **mixedcase**: `security` → `SeCuRiTy` (alternates case)
- **suffix**: `security` → `security789!` (adds random suffix)
- **prefix**: `security` → `42!security` (adds random prefix)  
- **insert**: `security` → `sec7ur!ity` (inserts random characters)
- **hybrid**: `security` → `S3cur1ty!42` (combines strategies)

## Examples

### Generate Different Types of Passwords

```bash
# Default password (14 characters, uppercase, lowercase, symbols)
$ passgen
� Your Password:
┌────────────────┐
│ U$DD$fico*Q,.Y │
└────────────────┘
📊 Length: 14 | Character types: Lowercase, Uppercase, Symbols | Strength: Very Strong 💪
🔒 Security info: 88.0 bits entropy, cracks in 4892016 years
💬 Someone's taking this security thing seriously! 🌟

# Secure password with all character types (gets the ice cold message!)
$ passgen --secure -l 20
� Your Password:
┌──────────────────────┐
│ 7hRQj<=1YS64M-ECL3iM │
└──────────────────────┘
📊 Length: 20 | Character types: Lowercase, Uppercase, Numbers, Symbols | Strength: Extremely Strong 🔥
🔒 Security info: 129.2 bits entropy, cracks in 1.2e+19 years
💬 Brr, that's ice cold security! Even hackers are shivering! 🥶

# Medium strength password (with sarcastic feedback)
$ passgen -l 8
🎯 Your Password:
┌──────────┐
│ rM!?M-k_ │
└──────────┘
📊 Length: 8 | Character types: Lowercase, Uppercase, Symbols | Strength: Medium ⚡
🔒 Security info: 50.3 bits entropy, cracks in 11.4 minutes
💡 Suggestions:
   • Consider using 12+ characters for better security
💬 Well, it's... adequate. I guess that's something! 👍
```

### Advanced Usage

```bash
# Generate multiple passwords
$ passgen -c 2 -l 10
� Password 1:
┌────────────┐
│ Cm=tnmnB#w │
└────────────┘
📊 Length: 10 | Character types: Lowercase, Uppercase, Symbols | Strength: Strong 💯
────────────────────────────────────────────────────────────
🎯 Password 2:
┌────────────┐
│ :<S<VZCalp │
└────────────┘
📊 Length: 10 | Character types: Lowercase, Uppercase, Symbols | Strength: Strong 💯

# Check version
$ passgen --version
passgen version v1.1.0
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

### v1.1.0 (Latest) 🎉

**New Features:**
- **🎯 Word-Based Password Generation**: Transform memorable words into secure passwords
  - 6 transformation strategies: leetspeak, mixedcase, suffix, prefix, insert, hybrid
  - 3 complexity levels: low, medium, high
  - Generate multiple variations of the same word
  - Intelligent pattern detection and security optimization

**Improvements:**
- **🏗️ Clean Architecture**: Refactored codebase with Domain-Driven Design principles
- **🎨 Enhanced UI**: More prominent password display with better visual formatting
- **🧪 Comprehensive Testing**: 21 test cases covering all architectural layers
- **🔧 Code Quality**: Fixed linting issues, updated deprecated APIs
- **⚡ CI/CD Optimization**: Enhanced GitHub Actions with better caching and security scanning

**Examples:**
```bash
# Transform words into secure passwords
passgen word "sunshine" --strategy hybrid --complexity high
# Output: Sunsh1n3!42

passgen word "coffee" --strategy leetspeak
# Output: c0ff33

passgen word "team" --count 3 --complexity medium
# Generate 3 variations: Team42, T3am78, Te@m91
```

**Breaking Changes:** None - fully backward compatible

### v1.0.0
- Initial release
- Basic password generation with customizable options
- Password strength checker
- Preset configurations
- Cross-platform support
