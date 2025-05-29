#!/bin/zsh

# Script to set up the development environment for the starknode-kit project

# --- Configuration ---
MIN_GO_VERSION="1.21" # Minimum recommended Go version

# --- Helper Functions ---
check_go_version() {
    if ! command -v go &> /dev/null; then
        echo "Go is not installed."
        echo "Please install Go version ${MIN_GO_VERSION} or higher."
        echo "Visit https://golang.org/doc/install for installation instructions."
        return 1
    fi

    current_go_version=$(go version | awk '{print $3}' | sed 's/go//')
    
    # Simple version comparison (adjust if more complex logic is needed)
    if printf "%s\n%s\n" "$MIN_GO_VERSION" "$current_go_version" | sort -V -C; then
        echo "Go version ${current_go_version} is installed."
    else
        echo "Go version ${current_go_version} is installed, but version ${MIN_GO_VERSION} or higher is recommended."
        echo "Please consider upgrading Go. Visit https://golang.org/doc/install for instructions."
        # Optionally, you could choose to exit here if the version is too old to proceed
        # return 1 
    fi
    return 0
}

# --- Main Setup ---
echo "Starting development environment setup for starknode-kit..."
echo ""

# 1. Check Go installation
echo "Step 1: Checking Go installation..."
if ! check_go_version; then
    echo ""
    echo "Please install or upgrade Go as instructed above and re-run this script."
    exit 1
fi
echo "Go installation check complete."
echo ""

# 2. Install Go module dependencies
echo "Step 2: Installing Go module dependencies..."
if go mod download; then
    echo "Go module dependencies downloaded successfully."
else
    echo "Failed to download Go module dependencies. Please check for errors."
    exit 1
fi
echo ""

# 3. Information about building and testing
echo "Step 3: Build and Test Information"
echo "The project uses a Makefile for common tasks:"
echo "  - To build the project: make build"
echo "  - To run tests:         make test"
echo "  - To clean artifacts:   make clean"
echo ""

echo "Development environment setup complete!"
echo "You should now be able to build and test the project."

exit 0
