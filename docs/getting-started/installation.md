# Installation and Setup

This guide will walk you through installing Starknode Kit on your system. We provide multiple installation methods to suit different preferences and environments.

## Prerequisites

Before installing Starknode Kit, ensure you have the following installed:

- **Go**: Version 1.24 or later ([Download](https://go.dev/dl/))
- **Rust**: For building Starknet clients ([Install Rust](https://rustup.rs/))
- **Make**: For building certain clients
  - Ubuntu/Debian: `sudo apt install make`
  - macOS: `brew install make`
  - Windows (WSL): `sudo apt install make`

## Installation Methods

### Method 1: Install Script (Recommended)

The easiest way to install Starknode Kit is using our installation script:

```bash
curl -sSL https://raw.githubusercontent.com/thebuidl-grid/starknode-kit/main/install.sh | bash
```

Or download and run manually:

```bash
wget https://raw.githubusercontent.com/thebuidl-grid/starknode-kit/main/install.sh
chmod +x install.sh
./install.sh
```

### Method 2: Go Install

If you have Go installed, you can install directly:

```bash
go install github.com/thebuidl-grid/starknode-kit@latest
```

### Method 3: Build from Source

For development or custom builds:

```bash
# Clone the repository
git clone https://github.com/thebuidl-grid/starknode-kit.git
cd starknode-kit

# Build and install
go build -o starknode-kit .
sudo mv starknode-kit /usr/local/bin/
```

## Verify Installation

After installation, verify that Starknode Kit is working correctly:

```bash
starknode-kit --help
```

You should see the help output with available commands.

## Initialize Configuration

Create your initial configuration file:

```bash
starknode-kit init
```

This creates a default configuration file in your home directory that you can customize for your needs.

## Next Steps

Now that you have Starknode Kit installed:

1. [Review Hardware Requirements](hardware-requirements.md)
2. [Follow the Quick Start Guide](quick-start.md)
3. [Configure your first client](configuration.md)

## Troubleshooting Installation

### Common Issues

**Permission Denied Error**
```bash
sudo chmod +x /usr/local/bin/starknode-kit
```

**Command Not Found**
Ensure `/usr/local/bin` is in your PATH:
```bash
echo 'export PATH=$PATH:/usr/local/bin' >> ~/.bashrc
source ~/.bashrc
```

**Go Version Issues**
Update Go to version 1.24 or later:
```bash
# Check current version
go version

# Update if needed
# Visit https://go.dev/dl/ for latest version
```

### Getting Help

If you encounter issues during installation:

- Check our [Troubleshooting Guide](../operations/troubleshooting.md)
- Join our [Telegram community](https://t.me/+SCPbza9fk8dkYWI0)
- Open an issue on [GitHub](https://github.com/thebuidl-grid/starknode-kit/issues)
