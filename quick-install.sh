#!/bin/bash

# Quick passgen installer - one-liner version
# Usage: curl -sSL https://raw.githubusercontent.com/kumarasakti/passgen/main/quick-install.sh | bash

set -e

echo "🚀 Quick installing passgen..."

# Install passgen
go install github.com/kumarasakti/passgen@latest

# Add to PATH for current session
export PATH=$PATH:$(go env GOPATH)/bin

# Add to shell config
SHELL_CONFIG=""
if [ -n "$ZSH_VERSION" ] || [[ "$SHELL" == *"zsh"* ]]; then
    SHELL_CONFIG="$HOME/.zshrc"
elif [ -n "$BASH_VERSION" ] || [[ "$SHELL" == *"bash"* ]]; then
    SHELL_CONFIG="$HOME/.bashrc"
else
    SHELL_CONFIG="$HOME/.profile"
fi

# Check if already in config
if ! grep -q "go env GOPATH" "$SHELL_CONFIG" 2>/dev/null; then
    echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> "$SHELL_CONFIG"
    echo "✅ Added Go bin to PATH in $SHELL_CONFIG"
fi

echo "🎉 passgen installed! Try: passgen --version" 