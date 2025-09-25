# Adding Clients

This guide covers how to add and configure different types of clients (Ethereum execution, consensus, and Starknet) to your Starknode Kit setup.

## Overview

Starknode Kit supports three types of clients:

- **Execution Clients**: Handle transaction execution and state management
- **Consensus Clients**: Handle block validation and consensus
- **Starknet Clients**: Handle Starknet L2 operations

## Adding Execution Clients

### Supported Execution Clients

| Client | Language | Description |
|--------|----------|-------------|
| **Geth** | Go | Most popular Ethereum client |
| **Reth** | Rust | High-performance Rust implementation |

### Add Geth

```bash
starknode-kit add --execution_client geth
```

### Add Reth

```bash
starknode-kit add --execution_client reth
```

### Execution Client Features

**Geth**:
- Mature and stable
- Extensive documentation
- Large community support
- Good for most use cases

**Reth**:
- High performance
- Modern architecture
- Fast sync capabilities
- Growing ecosystem

## Adding Consensus Clients

### Supported Consensus Clients

| Client | Language | Description |
|--------|----------|-------------|
| **Lighthouse** | Rust | High-performance consensus client |
| **Prysm** | Go | Feature-rich consensus client |

### Add Lighthouse

```bash
starknode-kit add --consensus_client lighthouse
```

### Add Prysm

```bash
starknode-kit add --consensus_client prysm
```

### Consensus Client Features

**Lighthouse**:
- Fast sync performance
- Low resource usage
- Strong security focus
- Good for production

**Prysm**:
- Rich feature set
- Web UI included
- Extensive tooling
- Good for development

## Adding Starknet Clients

### Supported Starknet Clients

| Client | Language | Description |
|--------|----------|-------------|
| **Juno** | Go | Full Starknet node implementation |

### Add Juno

```bash
starknode-kit add --starknet_client juno
```

### Juno Features

- Complete JSON-RPC support
- Fast synchronization
- Small database footprint
- WebSocket interface
- Production-ready

## Adding Client Pairs

### Complete Ethereum Setup

Add both execution and consensus clients together:

```bash
# Geth + Lighthouse
starknode-kit add --execution_client geth --consensus_client lighthouse

# Reth + Prysm
starknode-kit add --execution_client reth --consensus_client prysm

# Geth + Prysm
starknode-kit add --execution_client geth --consensus_client prysm
```

### Full Stack Setup

Add all three client types:

```bash
starknode-kit add --execution_client geth --consensus_client lighthouse --starknet_client juno
```

## Client Installation Process

When you add a client, Starknode Kit:

1. **Downloads** the client binary
2. **Installs** it to the appropriate location
3. **Configures** default settings
4. **Updates** your configuration file
5. **Validates** the installation

### Installation Locations

- **Linux/macOS**: `~/.starknode-kit/bin/`
- **Windows**: `%USERPROFILE%\.starknode-kit\bin\`

## Configuration After Adding

### View Added Clients

```bash
# Show all configured clients
starknode-kit config

# Show specific client type
starknode-kit config el    # Execution clients
starknode-kit config cl    # Consensus clients
starknode-kit config starknet  # Starknet clients
```

### Customize Client Settings

After adding, you can customize settings:

```bash
# Change ports
starknode-kit config el client=geth port=8545,8546

# Set execution type
starknode-kit config el client=geth execution_type=archive

# Add custom arguments
starknode-kit config el client=geth additional_args="--maxpeers,50,--cache,4096"
```

## Client Compatibility

### Recommended Combinations

**Production**:
- Geth + Lighthouse + Juno
- Reth + Lighthouse + Juno

**Development**:
- Geth + Prysm + Juno
- Reth + Prysm + Juno

**High Performance**:
- Reth + Lighthouse + Juno

### Network Compatibility

All clients support:
- Ethereum Mainnet
- Sepolia Testnet
- Holesky Testnet
- Custom networks

## Troubleshooting Client Addition

### Common Issues

**Download Fails**
```bash
# Check internet connection
ping google.com

# Check available disk space
df -h
```

**Permission Denied**
```bash
# Fix permissions
chmod +x ~/.starknode-kit/bin/*
```

**Client Already Exists**
```bash
# Remove existing client first
starknode-kit remove --execution_client geth

# Then add again
starknode-kit add --execution_client geth
```

### Verification

Verify client installation:

```bash
# Check if binary exists
ls -la ~/.starknode-kit/bin/

# Test client version
~/.starknode-kit/bin/geth version
~/.starknode-kit/bin/lighthouse --version
~/.starknode-kit/bin/juno --version
```

## Client Updates

### Check for Updates

```bash
starknode-kit update
```

### Update Specific Client

```bash
# Update all clients
starknode-kit update

# Update specific client type
starknode-kit update --execution_client geth
starknode-kit update --consensus_client lighthouse
starknode-kit update --starknet_client juno
```

## Removing Clients

### Remove Execution Client

```bash
starknode-kit remove --execution_client geth
```

### Remove Consensus Client

```bash
starknode-kit remove --consensus_client lighthouse
```

### Remove Starknet Client

```bash
starknode-kit remove --starknet_client juno
```

### Remove All Clients

```bash
starknode-kit remove --execution_client geth --consensus_client lighthouse --starknet_client juno
```

## Best Practices

### Client Selection

**For Beginners**:
- Start with Geth + Lighthouse
- Use default configurations
- Focus on stability

**For Production**:
- Use Reth + Lighthouse for performance
- Configure resource limits
- Set up monitoring

**For Development**:
- Use Geth + Prysm for features
- Enable debug logging
- Use testnets

### Configuration Management

- Keep configurations in version control
- Document custom settings
- Test changes on testnets first
- Monitor performance after changes

## Getting Help

For client-related issues:

- Check the [Troubleshooting Guide](../operations/troubleshooting.md)
- Review client-specific documentation
- Join our [Telegram community](https://t.me/+SCPbza9fk8dkYWI0)
- Open an issue on [GitHub](https://github.com/thebuidl-grid/starknode-kit/issues)
