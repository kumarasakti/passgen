# passgen

A secure, customizable password generator CLI tool written in Go with clean architecture.

```
    ____  ____ _______________ ____  ____ 
   / __ \/ __ `/ ___/ ___/ __ `/ _ \/ __ \
  / /_/ / /_/ (__  |__  ) /_/ /  __/ / / /
 / .___/\__,_/____/____/\__, /\___/_/ /_/ 
/_/                    /____/             

  passgen v1.1.0
  Cryptographically secure password generation
```

## Features

- **🔐 Secure Random Generation** — Uses `crypto/rand` for cryptographically secure randomness
- **🎨 Customizable Character Sets** — Lowercase, uppercase, numbers, symbols (toggle individually)
- **🔄 No-Repeat Mode** — `--no-repeat` flag guarantees no duplicate characters with full type coverage
- **🎯 Word-Based Passwords** — Transform memorable words into secure passwords (6 strategies, 3 complexity levels)
- **🔍 Password Strength Checker** — Analyze strength and get improvement suggestions
- **🚀 Preset Configurations** — Quick presets: secure, simple, pin, alphanumeric
- **📦 Batch Generation** — Generate multiple unique passwords at once
- **🌍 Cross-Platform** — Linux, macOS, Windows

## Installation

### Quick Install (Recommended)

```bash
curl -sSL https://raw.githubusercontent.com/kumarasakti/passgen/main/install.sh | bash
```

### Go Install

```bash
go install github.com/kumarasakti/passgen@latest
```

> After `go install`, add Go's bin to PATH: `export PATH=$PATH:$(go env GOPATH)/bin`

### PowerShell

```powershell
irm https://raw.githubusercontent.com/kumarasakti/passgen/main/install.ps1 | iex
```

### Build from Source

```bash
git clone https://github.com/kumarasakti/passgen.git
cd passgen
make build
sudo mv build/passgen /usr/local/bin/
```

### Pre-built Binaries

Download from the [releases page](https://github.com/kumarasakti/passgen/releases).

### Verify

```bash
passgen --version
passgen --help
```

## Usage

### Basic Generation

```bash
passgen                          # Default 14-char password (lower, upper, symbols)
passgen -l 16 -n                 # 16 chars with numbers
passgen --secure                 # All character types enabled
passgen -c 5                     # Generate 5 passwords
```

### Advanced Options

```bash
passgen --no-repeat -l 20 --secure          # No duplicate characters
passgen --exclude-similar                   # Exclude i, l, 1, L, o, 0, O
passgen --exclude "aeiou"                   # Exclude specific characters
passgen --lower=false --upper=false -n -l 6 # PIN (numbers only)
passgen --alphanumeric -l 12               # Letters and numbers only
```

### Word-Based Passwords

```bash
passgen word "security"                          # Default hybrid transformation
passgen word "sunshine" --strategy leetspeak    # s3cur1ty style
passgen word "secret" --complexity high          # Maximum complexity
passgen word "team" --count 3                    # 3 variations
```

**Strategies:** `leetspeak` · `mixedcase` · `suffix` · `prefix` · `insert` · `hybrid` (default)

**Complexity:** `low` · `medium` (default) · `high`

### Presets

```bash
passgen preset secure        # 16 chars, all types
passgen preset simple        # 12 chars, letters + numbers
passgen preset pin           # 6 digits
passgen preset alphanumeric  # 12 chars, letters + numbers
```

### Strength Checker

```bash
passgen check "mypassword123"
```

## Command Line Options

### Standard Generation

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--length` | `-l` | Password length | 14 |
| `--count` | `-c` | Number of passwords to generate | 1 |
| `--lower` | | Include lowercase letters | true |
| `--upper` | | Include uppercase letters | true |
| `--numbers` | `-n` | Include numbers | false |
| `--symbols` | `-s` | Include symbols | true |
| `--no-repeat` | | Avoid duplicate characters (guaranteed type coverage) | false |
| `--exclude-similar` | | Exclude similar characters (il1Lo0O) | false |
| `--exclude` | | Characters to exclude | "" |
| `--secure` | `-S` | Enable all character types | false |
| `--simple` | `-m` | Letters + numbers only | false |
| `--alphanumeric` | `-a` | Alphanumeric only | false |
| `--help` | `-h` | Show help | |
| `--version` | `-v` | Show version | |

### Word-Based Generation

```bash
passgen word <word> [flags]
```

| Flag | Short | Description | Default | Options |
|------|-------|-------------|---------|---------|
| `--strategy` | `-s` | Transformation strategy | hybrid | leetspeak, mixedcase, suffix, prefix, insert, hybrid |
| `--complexity` | `-x` | Complexity level | medium | low, medium, high |
| `--count` | `-c` | Number of variations | 1 | |

## Examples

```bash
$ passgen
🎯 Your Password:
┌────────────────┐
│ U$DD$fico*Q,.Y │
└────────────────┘
📊 Length: 14 | Character types: Lowercase, Uppercase, Symbols | Strength: Very Strong 💪
🔒 Security info: 88.0 bits entropy, cracks in 4892016 years
💬 Someone's taking this security thing seriously! 🌟

$ passgen --secure -l 20 --no-repeat
🎯 Your Password:
┌──────────────────────┐
│ dP$<c1Xk6.q9!jLv7bRz │
└──────────────────────┘
📊 Length: 20 | Character types: Lowercase, Uppercase, Numbers, Symbols | Strength: Extremely Strong 🔥
🔒 Security info: 101.3 bits entropy, cracks in 5.4e+19 years
💬 Brr, that's ice cold security! Even hackers are shivering! 🥶
```

## Development

```bash
make deps       # Install dependencies
make build      # Build for current platform
make build-all  # Build for all platforms
make test       # Run all tests
make dev        # Format, lint, test, build
make release    # Create release archives
```

## Security

This tool uses Go's `crypto/rand` package for cryptographically secure random number generation. Generated passwords are suitable for production use.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License — see [LICENSE](LICENSE) for details.

## Changelog

### v1.1.0

- **Word-Based Password Generation** — 6 transformation strategies, 3 complexity levels
- **No-Repeat Mode** — `--no-repeat` flag for no-duplicate passwords with guaranteed type coverage
- **Batch Uniqueness** — `GenerateMultiplePasswords` retries on collisions
- **BuildCategories** — Per-category charset support for guaranteed character type coverage
- **Fisher-Yates Shuffle** — Cryptographically secure shuffle for no-repeat mode
- **Clean Architecture** — Domain-Driven Design with domain/application/infrastructure layers
- **Comprehensive Testing** — 30+ test cases across all architectural layers
- Fully backward compatible

### v1.0.0

- Initial release
- Basic password generation with customizable options
- Password strength checker
- Preset configurations
- Cross-platform support
