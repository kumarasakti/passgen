#!/bin/bash

# passgen Installation Script
# This script installs passgen and automatically configures your PATH

set -e

echo "🚀 Installing passgen..."

# Install using go install
go install github.com/kumarasakti/passgen@latest

# Get Go bin path
GOBIN=$(go env GOPATH)/bin
PASSGEN_PATH="$GOBIN/passgen"

# Check if installation was successful
if [ ! -f "$PASSGEN_PATH" ]; then
    echo "❌ Installation failed. Please check your Go installation."
    exit 1
fi

echo "✅ passgen binary installed to: $PASSGEN_PATH"

# Function to detect user's shell
detect_shell() {
    # Check current shell from SHELL environment variable first
    case "$SHELL" in
        */zsh) echo "zsh" ;;
        */bash) echo "bash" ;;
        */fish) echo "fish" ;;
        *) 
            # Fallback to checking shell variables
            if [ -n "$ZSH_VERSION" ]; then
                echo "zsh"
            elif [ -n "$BASH_VERSION" ]; then
                echo "bash"
            else
                echo "unknown"
            fi
            ;;
    esac
}

# Function to get shell config file
get_shell_config() {
    local shell_type=$1
    case $shell_type in
        zsh)
            echo "$HOME/.zshrc"
            ;;
        bash)
            if [ -f "$HOME/.bashrc" ]; then
                echo "$HOME/.bashrc"
            elif [ -f "$HOME/.bash_profile" ]; then
                echo "$HOME/.bash_profile"
            else
                echo "$HOME/.bashrc"
            fi
            ;;
        fish)
            echo "$HOME/.config/fish/config.fish"
            ;;
        *)
            echo "$HOME/.profile"
            ;;
    esac
}

# Function to check if Go bin is already in PATH
is_in_path() {
    case ":$PATH:" in
        *":$GOBIN:"*) return 0 ;;
        *) return 1 ;;
    esac
}

# Main PATH configuration logic
configure_path() {
    if is_in_path; then
        echo "✅ Go bin directory is already in your PATH"
        return 0
    fi

    echo "🔧 Configuring PATH..."
    
    local shell_type=$(detect_shell)
    local config_file=$(get_shell_config $shell_type)
    local path_export='export PATH=$PATH:$(go env GOPATH)/bin'
    
    echo "Detected shell: $shell_type"
    echo "Config file: $config_file"
    
    # Check if the PATH export already exists in the config file
    if [ -f "$config_file" ] && grep -q "go env GOPATH" "$config_file"; then
        echo "⚠️  Go PATH configuration already exists in $config_file"
    else
        # Create config file directory if it doesn't exist (for fish)
        mkdir -p "$(dirname "$config_file")"
        
        # Add PATH export to config file
        echo "" >> "$config_file"
        echo "# Added by passgen installer" >> "$config_file"
        echo "$path_export" >> "$config_file"
        echo "✅ Added Go bin to PATH in $config_file"
    fi
    
    # Export PATH for current session
    export PATH=$PATH:$GOBIN
    
    return 0
}

# Configure PATH
configure_path

# Test if passgen is now accessible
if command -v passgen >/dev/null 2>&1; then
    echo ""
    echo "🎉 Installation successful!"
    echo "✅ passgen is now available in your PATH"
    echo ""
    echo "Try it out:"
    echo "  passgen --version"
    echo "  passgen --help"
    echo "  passgen"
    echo ""
else
    echo ""
    echo "⚠️  Installation completed but passgen is not immediately available."
    echo "Please restart your terminal or run:"
    echo "  source $(get_shell_config $(detect_shell))"
    echo ""
    echo "Or run passgen with full path:"
    echo "  $PASSGEN_PATH"
    echo ""
fi

echo "📖 For more information, visit: https://github.com/kumarasakti/passgen" 