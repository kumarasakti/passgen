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

### Method 2: Quick One-liner

```bash
curl -sSL https://raw.githubusercontent.com/kumarasakti/passgen/main/quick-install.sh | bash
```

### Method 3: Manual Go Install

```bash
go install github.com/kumarasakti/passgen@latest
```

**Important**: After manual installation, you need to add Go's bin directory to your PATH:

```bash
# For zsh users (default on macOS)
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
source ~/.zshrc

# For bash users
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc
```

### Method 4: Download Pre-built Binaries

Download the latest release for your platform from the [releases page](https://github.com/kumarasakti/passgen/releases) and place it in your PATH.

### Method 5: Build from Source

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
# Default password
$ passgen
Kj8mN2pL9xQr

# Secure password with all character types
$ passgen --secure -l 20
Kj8mN2pL9xQr!@#$%^&*

# Simple password (no symbols)
$ passgen --simple -l 16
Kj8mN2pL9xQrTyUi

# PIN number
$ passgen preset pin
582947

# Multiple passwords
$ passgen -c 3 -l 10
Kj8mN2pL9x
Qr3tY6uI8o
Pl4sD7fG9h
```

### Advanced Usage

```bash
# Exclude confusing characters
$ passgen --exclude-similar -l 12
KjmnpqrTyUiE

# Custom character exclusion
$ passgen --exclude "aeiou" -l 12
KjmnpqrTyUiE

# Only uppercase and numbers
$ passgen --lower=false --symbols=false -l 8
KJ8MN2PL
```

## Development

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
