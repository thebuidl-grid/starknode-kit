#!/bin/bash

# Configuration - Update these variables for your specific package
GITHUB_REPO="thebuidl-grid/starknode-kit" 
BINARY_NAME="starknode-kit" 
INSTALL_DIR="/usr/local/bin" 

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
print_status "Checking prerequisites..."

if ! command_exists git; then
    print_error "git is required but not installed. Please install git first."
    exit 1
fi

if ! command_exists go; then
    print_error "Go is required but not installed. Please install Go first."
    exit 1
fi

# Create temporary directory
TEMP_DIR=$(mktemp -d)
print_status "Created temporary directory: $TEMP_DIR"

# Cleanup function
cleanup() {
    print_status "Cleaning up temporary files..."
    rm -rf "$TEMP_DIR"
}

# Set trap to cleanup on exit
trap cleanup EXIT

# Clone the repository
print_status "Cloning repository: https://github.com/$GITHUB_REPO"
if ! git clone "https://github.com/$GITHUB_REPO.git" "$TEMP_DIR/$BINARY_NAME"; then
    print_error "Failed to clone repository"
    exit 1
fi

# Change to project directory
cd "$TEMP_DIR/$BINARY_NAME" || {
    print_error "Failed to change to project directory"
    exit 1
}

# Check if go.mod exists (Go modules)
if [ -f "go.mod" ]; then
    print_status "Go modules detected, downloading dependencies..."
    go mod download
else
    print_status "No go.mod found, assuming GOPATH mode..."
fi

# Build the application
print_status "Building the application..."
if ! go build -o "$BINARY_NAME" .; then
    print_error "Failed to build the application"
    exit 1
fi

# Check if binary was created
if [ ! -f "$BINARY_NAME" ]; then
    print_error "Binary was not created successfully"
    exit 1
fi

# Create install directory if it doesn't exist
if [ ! -d "$INSTALL_DIR" ]; then
    print_warning "Install directory $INSTALL_DIR does not exist, creating it..."
    sudo mkdir -p "$INSTALL_DIR"
fi

# Install the binary
print_status "Installing $BINARY_NAME to $INSTALL_DIR..."
if ! sudo cp "$BINARY_NAME" "$INSTALL_DIR/"; then
    print_error "Failed to install binary to $INSTALL_DIR"
    exit 1
fi

# Make it executable
sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"

# Verify installation
if command_exists "$BINARY_NAME"; then
    print_status "✓ Installation successful!"
    print_status "You can now use '$BINARY_NAME' from anywhere in your terminal"
    
    # Show version if available
    if "$BINARY_NAME" --version >/dev/null 2>&1; then
        VERSION=$("$BINARY_NAME" --version)
        print_status "Installed version: $VERSION"
    elif "$BINARY_NAME" -version >/dev/null 2>&1; then
        VERSION=$("$BINARY_NAME" -version)
        print_status "Installed version: $VERSION"
    fi
else
    print_warning "Installation completed but '$BINARY_NAME' is not in PATH"
    print_warning "You may need to add $INSTALL_DIR to your PATH or restart your terminal"
fi

print_status "Installation complete!"
