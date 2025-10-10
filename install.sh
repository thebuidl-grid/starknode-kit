#!/bin/bash

# Configuration - Update these variables for your specific package
GITHUB_REPO="thebuidl-grid/starknode-kit"
BINARY_NAME="starknode-kit"
INSTALL_DIR="/usr/local/bin"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Node type selection variables
SELECTED_NODE_TYPE=""
SELECTED_NETWORK=""
SELECTED_EL_CLIENT=""
SELECTED_CL_CLIENT=""

IS_STARKNET_NODE=1
IS_VALIDATOR_NODE=1

handle_keyboard_interrupt() {
    clear
    print_status "Keyboard Interrupt"
    exit 1
}

# Function to display figlet banner
show_banner() {
    clear
    if command_exists figlet; then
        figlet -c "StarkNode Kit"
        echo
        echo -e "${CYAN}════════════════════════════════════════════════════════════════════${NC}"
        echo -e "${CYAN}              Welcome to the StarkNode Kit Installer                ${NC}"
        echo -e "${CYAN}════════════════════════════════════════════════════════════════════${NC}"
    else
        echo -e "${CYAN}╔══════════════════════════════════════════════════════════════════╗${NC}"
        echo -e "${CYAN}║                                                                  ║${NC}"
        echo -e "${CYAN}║              ███████╗████████╗ █████╗ ██████╗ ██╗  ██╗           ║${NC}"
        echo -e "${CYAN}║              ██╔════╝╚══██╔══╝██╔══██╗██╔══██╗██║ ██╔╝           ║${NC}"
        echo -e "${CYAN}║              ███████╗   ██║   ███████║██████╔╝█████╔╝            ║${NC}"
        echo -e "${CYAN}║              ╚════██║   ██║   ██╔══██║██╔══██╗██╔═██╗            ║${NC}"
        echo -e "${CYAN}║              ███████║   ██║   ██║  ██║██║  ██║██║  ██╗           ║${NC}"
        echo -e "${CYAN}║              ╚══════╝   ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝           ║${NC}"
        echo -e "${CYAN}║                                                                  ║${NC}"
        echo -e "${CYAN}║                         NODE KIT INSTALLER                       ║${NC}"
        echo -e "${CYAN}║                                                                  ║${NC}"
        echo -e "${CYAN}╚══════════════════════════════════════════════════════════════════╝${NC}"
    fi
    echo
    echo -e "${GREEN}        Build, Deploy, and Manage Starknet Infrastructure${NC}"
    echo
}

# Function to show node selection menu
show_node_selection() {
    while true; do
        clear
        show_banner
        
        echo -e "${BLUE}Please select the type of node you want to set up:${NC}"
        echo
        echo -e "${YELLOW}1)${NC} Ethereum Full Node"
        echo -e "   ${CYAN}└── Run a complete Ethereum node with full blockchain data${NC}"
        echo
        echo -e "${YELLOW}2)${NC} Starknet Full Node"
        echo -e "   ${CYAN}└── Run a complete Starknet node for Layer 2 scaling${NC}"
        echo
        echo -e "${YELLOW}3)${NC} Starknet Validator Node"
        echo -e "   ${CYAN}└── Participate in Starknet consensus and earn rewards${NC}"
        echo
        echo -e "${YELLOW}4)${NC} Update Starknod-kit"
        echo -e "   ${CYAN}└── Update starknode-kit${NC}"
        echo
        echo -e "${RED}5)${NC} Exit"
        echo
        echo -n -e "${GREEN}Enter your choice [1-5]: ${NC}"
        
        read -r choice
        
        case $choice in
            1)
                SELECTED_NODE_TYPE="ethereum"
                print_status "Selected: Ethereum Full Node"
                break
                ;;
            2)
                SELECTED_NODE_TYPE="starknet"
                print_status "Selected: Starknet Full Node"
                break
                ;;
            3)
                SELECTED_NODE_TYPE="validator"
                print_status "Selected: Starknet Validator Node"
                break
                ;;
            4)
                SELECTED_NODE_TYPE="update"
                print_status "Update client"
                break
                ;;
            5)
                clear
                print_status "Installation cancelled by user"
                exit 0
                ;;
            *)
                print_error "Invalid choice. Please select 1-5."
                sleep 2
                ;;
        esac
    done
    
    echo
    sleep 1
}

# Function to select Ethereum network
select_ethereum_network() {
    while true; do
        clear
        show_banner
        
        echo -e "${BLUE}Select Ethereum Network:${NC}"
        echo
        echo -e "${YELLOW}1)${NC} Mainnet"
        echo -e "   ${CYAN}└── Ethereum main network (production)${NC}"
        echo
        echo -e "${YELLOW}2)${NC} Sepolia"
        echo -e "   ${CYAN}└── Ethereum test network${NC}"
        echo
        echo -n -e "${GREEN}Enter your choice [1-2]: ${NC}"
        
        read -r choice
        
        case $choice in
            1)
                SELECTED_NETWORK="mainnet"
                print_status "Selected: Ethereum Mainnet"
                break
                ;;
            2)
                SELECTED_NETWORK="sepolia"
                print_status "Selected: Sepolia Testnet"
                break
                ;;
            *)
                print_error "Invalid choice. Please select 1-2."
                sleep 2
                ;;
        esac
    done
    echo
    sleep 1
}

# Function to select Execution Layer client
select_el_client() {
    while true; do
        clear
        show_banner
        
        echo -e "${BLUE}Select Execution Layer (EL) Client:${NC}"
        echo
        echo -e "${YELLOW}1)${NC} Geth"
        echo -e "   ${CYAN}└── Go Ethereum (Most popular, stable)${NC}"
        echo
        echo -e "${YELLOW}2)${NC} Reth"
        echo -e "   ${CYAN}└── Rust Ethereum (High performance, modern)${NC}"
        echo
        echo -n -e "${GREEN}Enter your choice [1-2]: ${NC}"
        
        read -r choice
        
        case $choice in
            1)
                SELECTED_EL_CLIENT="geth"
                print_status "Selected: Geth (Go Ethereum)"
                break
                ;;
            2)
                SELECTED_EL_CLIENT="reth"
                print_status "Selected: Reth (Rust Ethereum)"
                break
                ;;
            *)
                print_error "Invalid choice. Please select 1-2."
                sleep 2
                ;;
        esac
    done
    echo
    sleep 1
}

# Function to select Consensus Layer client
select_cl_client() {
    while true; do
        clear
        show_banner
        
        echo -e "${BLUE}Select Consensus Layer (CL) Client:${NC}"
        echo
        echo -e "${YELLOW}1)${NC} Lighthouse"
        echo -e "   ${CYAN}└── Rust implementation (Efficient, reliable)${NC}"
        echo
        echo -e "${YELLOW}2)${NC} Prysm"
        echo -e "   ${CYAN}└── Go implementation (Feature-rich, popular)${NC}"
        echo
        echo -n -e "${GREEN}Enter your choice [1-2]: ${NC}"
        
        read -r choice
        
        case $choice in
            1)
                SELECTED_CL_CLIENT="lighthouse"
                print_status "Selected: Lighthouse"
                break
                ;;
            2)
                SELECTED_CL_CLIENT="prysm"
                print_status "Selected: Prysm"
                break
                ;;
            *)
                print_error "Invalid choice. Please select 1-2."
                sleep 2
                ;;
        esac
    done
    echo
    sleep 1
}

# Function to handle complete Ethereum selection flow
handle_ethereum_selection() {
    select_ethereum_network
    select_el_client
    select_cl_client
    show_node_config
}

show_node_config() {
    echo -e "${BLUE}Ethereum Node Configuration Summary:${NC}"
    echo -e "${CYAN}Network:${NC} $SELECTED_NETWORK"
    echo -e "${CYAN}Execution Client:${NC} $SELECTED_EL_CLIENT"
    echo -e "${CYAN}Consensus Client:${NC} $SELECTED_CL_CLIENT"
    if [ "$IS_STARKNET_NODE" != 1 ]; then
      echo -e "${CYAN}Starknet Client:${NC} Juno"
    fi
    echo
    echo -e "${YELLOW}Requirements:${NC}"
    if [ "$SELECTED_NETWORK" == "mainnet" ]; then
        echo "• Disk Space: 1TB+ (mainnet)"
        echo "• RAM: 16GB+ recommended"
        echo "• Sync Time: 1-3 days"
    else
        echo "• Disk Space: 100GB+ (sepolia)"
        echo "• RAM: 8GB+ recommended"
        echo "• Sync Time: 2-6 hours"
    fi
    echo
}

# Function to display node type specific information
show_node_info() {
    case $SELECTED_NODE_TYPE in
        "ethereum")
            echo -e "${BLUE}Ethereum Full Node Setup:${NC}"
            echo "• Downloads and syncs the complete Ethereum blockchain"
            echo "• Requires significant disk space (1TB+ recommended)"
            echo "• Provides full validation and RPC capabilities"
            echo "• Estimated sync time: 1-3 days depending on hardware"
            ;;
        "starknet")
            echo -e "${BLUE}Starknet Full Node Setup:${NC}"
            echo "• Connects to Starknet Layer 2 network"
            echo "• Requires moderate disk space (100GB+ recommended)"
            echo "• Provides fast transaction processing"
            echo "• Estimated sync time: 2-6 hours"
            IS_STARKNET_NODE=0
            ;;
        "validator")
            echo -e "${BLUE}Starknet Validator Node Setup:${NC}"
            echo "• Participates in network consensus"
            echo "• Requires staking STRK tokens"
            echo "• Earns validation rewards"
            echo "• Requires high uptime and reliable connection"
            IS_STARKNET_NODE=0
            IS_VALIDATOR_NODE=0
            ;;
    esac
    echo
}

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

# Function to install figlet if not present
check_figlet() {
    if ! command_exists figlet; then
        print_warning "figlet not found. Attempting to install..."
        
        # Detect OS and install figlet
        if command_exists apt-get; then
            sudo apt-get update && sudo apt-get install -y figlet
        elif command_exists yum; then
            sudo yum install -y figlet
        elif command_exists brew; then
            brew install figlet
        elif command_exists pacman; then
            sudo pacman -S figlet
        else
            print_warning "Could not install figlet automatically. Banner will use ASCII art."
        fi
    fi
}

# Check prerequisites
check_prerequisites() {
    print_status "Checking prerequisites..."
    
    # Check for figlet and try to install if missing
    check_figlet
    
    if ! command_exists git; then
        print_error "git is required but not installed. Please install git first."
        exit 1
    fi
    
    if ! command_exists go; then
        print_error "Go is required but not installed. Please install Go first."
        exit 1
    fi
    
    print_status "All prerequisites satisfied!"
    echo
}

# Main installation process
perform_installation() {
    if [ "$USE_LOCAL" = true ]; then
        print_status "Using local main.go file for installation. No binary will be installed."
    
        # Check if go.mod exists (Go modules)
        if [ -f "go.mod" ]; then
            print_status "Go modules detected, downloading dependencies..."
            go mod download
        else
            print_status "No go.mod found, assuming GOPATH mode..."
        fi
    
        # Build the application with node type flag
        print_status "Preparing to configure clients using local source for $SELECTED_NODE_TYPE node..."
        BUILD_FLAGS=""
        case $SELECTED_NODE_TYPE in
            "ethereum") BUILD_FLAGS="-tags ethereum" ;;
            "starknet") BUILD_FLAGS="-tags starknet" ;;
            "validator") BUILD_FLAGS="-tags validator" ;;
        esac

        CONFIG_CMD="go run main.go config new -i"
        CONFIG_CMD="$CONFIG_CMD --network $SELECTED_NETWORK"
        CONFIG_CMD="$CONFIG_CMD --consensus-client $SELECTED_CL_CLIENT"
        CONFIG_CMD="$CONFIG_CMD --execution-client $SELECTED_EL_CLIENT"
        if [ "$IS_STARKNET_NODE" == 0 ]; then
            CONFIG_CMD="$CONFIG_CMD --starknet-node juno"
        fi
        if [ "$IS_VALIDATOR_NODE" == 0 ]; then
            CONFIG_CMD="$CONFIG_CMD --validator"
        fi

        print_status "Running configuration command: $CONFIG_CMD"
        if ! eval "$CONFIG_CMD"; then
            print_error "Failed to configure and install clients."
            exit 1
        fi
        print_status "Client configuration and installation complete!"

    else
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
    
        # Get version from git tag
        VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "dev")
        print_status "Building version: $VERSION"

        # Construct linker flags
        LDFLAGS="-X github.com/thebuidl-grid/starknode-kit/pkg/versions.StarkNodeVersion=$VERSION"

        # Build the application with node type flag
        print_status "Building the application for $SELECTED_NODE_TYPE node..."
        BUILD_FLAGS=""
        case $SELECTED_NODE_TYPE in
            "ethereum") BUILD_FLAGS="-tags ethereum" ;;
            "starknet") BUILD_FLAGS="-tags starknet" ;;
            "validator") BUILD_FLAGS="-tags validator" ;;
        esac
    
        if ! go build $BUILD_FLAGS -ldflags="$LDFLAGS" -o "$BINARY_NAME" .; then
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
    
        # Create config file with node type
        CONFIG_DIR="$HOME/.starknode-kit"
        mkdir -p "$CONFIG_DIR"
        echo "node_type=$SELECTED_NODE_TYPE" > "$CONFIG_DIR/config"
    
        # Verify installation
        if command_exists "$BINARY_NAME"; then
            echo
            print_status "✓ Installation successful!"
            print_status "You can now use '$BINARY_NAME' from anywhere in your terminal"
            print_status "Node type configured: $SELECTED_NODE_TYPE"
            
            # Show version if available
            if "$BINARY_NAME" --version >/dev/null 2>&1; then
                VERSION=$("$BINARY_NAME" --version)
                print_status "Installed version: $VERSION"
            elif "$BINARY_NAME" -version >/dev/null 2>&1; then
                VERSION=$("$BINARY_NAME" -version)
                print_status "Installed version: $VERSION"
            fi

            print_status "Configuring and installing clients based on your selections..."

            CONFIG_CMD="$BINARY_NAME config new -i"

            # Add network flag
            CONFIG_CMD="$CONFIG_CMD --network $SELECTED_NETWORK"

            CONFIG_CMD="$CONFIG_CMD --consensus-client $SELECTED_CL_CLIENT"
            CONFIG_CMD="$CONFIG_CMD --execution-client $SELECTED_EL_CLIENT"

            # Add starknet-node flag
            if [ "$IS_STARKNET_NODE" == 0 ]; then # 0 means it's a starknet node (starknet or validator)
                CONFIG_CMD="$CONFIG_CMD --starknet-node juno"
            fi

            # Add validator flag
            if [ "$IS_VALIDATOR_NODE" == 0 ]; then # 0 means it's a validator node
                CONFIG_CMD="$CONFIG_CMD --validator"
            fi

            print_status "Running configuration command: $CONFIG_CMD"
            if ! $CONFIG_CMD; then
                print_error "Failed to configure and install clients."
                exit 1
            fi
            print_status "Client configuration and installation complete!"

        else
            print_warning "Installation completed but '$BINARY_NAME' is not in PATH"
            print_warning "You may need to add $INSTALL_DIR to your PATH or restart your terminal"
        fi
    fi
}

# Function to show completion message
show_completion() {
    echo
    echo -e "${GREEN}════════════════════════════════════════════════════════════════════${NC}"
    if [ "$USE_LOCAL" = true ]; then
        echo -e "${GREEN}                      Configuration Complete!                      ${NC}"
    else
        echo -e "${GREEN}                        Installation Complete!                       ${NC}"
    fi
    echo -e "${GREEN}════════════════════════════════════════════════════════════════════${NC}"
    echo
    echo -e "${CYAN}Next Steps:${NC}"
    if [ "$USE_LOCAL" = true ]; then
        echo "1. To see available commands, run: go run main.go --help"
        echo "2. To initialize your $SELECTED_NODE_TYPE node, run: go run main.go init"
        echo "3. To start your node, run: go run main.go start"
    else
        echo "1. Run '$BINARY_NAME --help' to see available commands"
        echo "2. Initialize your $SELECTED_NODE_TYPE node with '$BINARY_NAME init'"
        echo "3. Start your node with '$BINARY_NAME start'"
    fi
    echo
    echo -e "${YELLOW}For support and documentation, visit:${NC}"
    echo "https://github.com/$GITHUB_REPO"
    echo
}

# Function to update the package
perform_update() {
    print_status "Checking for updates..."

    # Check for current installation
    if ! command_exists "$BINARY_NAME"; then
        print_error "$BINARY_NAME is not installed. Please run the full installation first."
        exit 1
    fi

    # Get current version
    CURRENT_VERSION_FULL=$($BINARY_NAME version 2>/dev/null)
    CURRENT_VERSION=$(echo "$CURRENT_VERSION_FULL" | grep -o -E 'v?[0-9]+\.[0-9]+\.[0-9]+' | head -n1)
    if [ -z "$CURRENT_VERSION" ]; then
        print_warning "Could not determine current version from output: '$CURRENT_VERSION_FULL'. Proceeding with update."
    else
        print_status "Current version: $CURRENT_VERSION"
    fi

    # Get latest version from GitHub using git tags
    if ! command_exists git; then
        print_warning "git is not installed, cannot check for latest version. Proceeding with update."
    else
        LATEST_VERSION=$(git ls-remote --tags "https://github.com/$GITHUB_REPO.git" | awk '{print $2}' | grep 'refs/tags/v' | sed 's|refs/tags/||' | sort -V | tail -n 1)
        if [ -z "$LATEST_VERSION" ]; then
            print_warning "Could not determine latest version from git tags. Proceeding with update."
        else
            print_status "Latest version available: $LATEST_VERSION"
        fi
    fi

    # Compare versions
    if [ -n "$CURRENT_VERSION" ] && [ -n "$LATEST_VERSION" ]; then
        # Strip 'v' prefix for comparison
        if [ "${CURRENT_VERSION#v}" == "${LATEST_VERSION#v}" ]; then
            print_status "You already have the latest version ($CURRENT_VERSION)."
            exit 0
        fi
        print_status "New version available. Proceeding with update..."
    fi

    print_status "Starting update process..."

    # Read node type from config
    CONFIG_FILE="$HOME/.config/starknode-kit/config"
    if [ -f "$CONFIG_FILE" ]; then
        # shellcheck source=/dev/null
        . "$CONFIG_FILE"
        print_status "Existing configuration found. Node type: $node_type"
        SELECTED_NODE_TYPE=$node_type
    else
        print_error "No existing configuration found at $CONFIG_FILE."
        print_error "Cannot determine node type for update. Please run the full installation first."
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

    # Get version from git tag
    VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "dev")
    print_status "Building version: $VERSION"

    # Construct linker flags
    LDFLAGS="-X github.com/thebuidl-grid/starknode-kit/pkg/versions.StarkNodeVersion=$VERSION"

    # Build the application with node type flag
    print_status "Building the application for $SELECTED_NODE_TYPE node..."
    BUILD_FLAGS=""
    case $SELECTED_NODE_TYPE in
        "ethereum") BUILD_FLAGS="-tags ethereum" ;;
        "starknet") BUILD_FLAGS="-tags starknet" ;;
        "validator") BUILD_FLAGS="-tags validator" ;;
    esac

    if ! go build $BUILD_FLAGS -ldflags="$LDFLAGS" -o "$BINARY_NAME" .; then
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
        echo
        print_status "✓ Update successful!"
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
        print_warning "Update completed but '$BINARY_NAME' is not in PATH"
        print_warning "You may need to add $INSTALL_DIR to your PATH or restart your terminal"
    fi
}

# Function to show update completion message
show_update_completion() {
    echo
    echo -e "${GREEN}════════════════════════════════════════════════════════════════════${NC}"
    echo -e "${GREEN}                         Update Complete!                          ${NC}"
    echo -e "${GREEN}════════════════════════════════════════════════════════════════════${NC}"
    echo
    echo -e "${CYAN}Next Steps:${NC}"
    echo "1. Run '$BINARY_NAME --help' to see available commands"
    echo "2. To start your node, run: '$BINARY_NAME start'"
    echo
    echo -e "${YELLOW}For support and documentation, visit:${NC}"
    echo "https://github.com/$GITHUB_REPO"
    echo
}

# Main execution
main() {
    USE_LOCAL=false
    for arg in "$@"; do
        if [ "$arg" == "--use-local" ]; then
            USE_LOCAL=true
        fi
    done

    trap handle_keyboard_interrupt SIGINT
    show_banner
    check_prerequisites
    show_node_selection

    if [ "$SELECTED_NODE_TYPE" == "update" ]; then
        perform_update
        show_update_completion
        exit 0
    fi

    handle_ethereum_selection
    case $SELECTED_NODE_TYPE in
        "starknet")
            IS_STARKNET_NODE=0
            ;;
        "validator")
            IS_STARKNET_NODE=0
            IS_VALIDATOR_NODE=0
            ;;
    esac

    perform_installation
    show_completion
}

# Run main function
main "$@" "$@"
