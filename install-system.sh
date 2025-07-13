#!/bin/bash

# passgen System-Wide Installation Script
# This script installs passgen to /usr/local/bin (requires sudo)

set -e

echo "🚀 Installing passgen system-wide..."

# Check if Go is installed
if ! command -v go >/dev/null 2>&1; then
    echo "❌ Go is not installed. Please install Go first: https://golang.org/dl/"
    exit 1
fi

# Create temporary directory
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

echo "📦 Downloading and building passgen..."

# Download and build
go mod init temp-passgen-install
go get github.com/kumarasakti/passgen@latest
go build -o passgen github.com/kumarasakti/passgen

# Check if build was successful
if [ ! -f "passgen" ]; then
    echo "❌ Build failed"
    exit 1
fi

echo "✅ Build successful"

# Install to system directory
echo "🔧 Installing to /usr/local/bin (requires sudo)..."

# Check if /usr/local/bin exists, create if not
if [ ! -d "/usr/local/bin" ]; then
    echo "Creating /usr/local/bin directory..."
    sudo mkdir -p /usr/local/bin
fi

# Install the binary
sudo cp passgen /usr/local/bin/passgen
sudo chmod +x /usr/local/bin/passgen

# Cleanup
cd /
rm -rf "$TEMP_DIR"

# Verify installation
if command -v passgen >/dev/null 2>&1; then
    echo ""
    echo "🎉 Installation successful!"
    echo "✅ passgen is now available system-wide"
    echo ""
    echo "Try it out:"
    echo "  passgen --version"
    echo "  passgen --help"
    echo "  passgen"
    echo ""
    echo "📍 Installed to: /usr/local/bin/passgen"
else
    echo ""
    echo "⚠️  Installation may have failed. Please check if /usr/local/bin is in your PATH"
    echo "You can run: echo \$PATH | grep -q /usr/local/bin && echo 'OK' || echo 'Not in PATH'"
fi

echo "📖 For more information, visit: https://github.com/kumarasakti/passgen" 