# Configuration

Starknode Kit uses a YAML configuration file to manage your node setup. This guide covers all configuration options and how to customize them for your needs.

## Configuration File Location

The configuration file is located at:
- **Linux/macOS**: `~/.starknode-kit/config.yaml`
- **Windows**: `%USERPROFILE%\.starknode-kit\config.yaml`

## Configuration Structure

```yaml
network: mainnet  # Network to connect to
execution_client:
  name: geth
  port: [8545, 8546]
  execution_type: full
consensus_client:
  name: lighthouse
  port: [9000, 9001]
  consensus_checkpoint: ""
starknet_client:
  name: juno
  port: [6060, 6061]
  eth_node_url: "http://localhost:8545"
```

## Network Configuration

### Supported Networks

| Network | Description | Use Case |
|---------|-------------|----------|
| `mainnet` | Ethereum mainnet | Production |
| `sepolia` | Ethereum testnet | Testing |
| `goerli` | Ethereum testnet (deprecated) | Legacy testing |
| `holesky` | Ethereum testnet | Testing |

### Setting Network

```bash
# Set network via command line
starknode-kit config -n sepolia

# Or edit config file directly
network: sepolia
```

## Execution Client Configuration

### Supported Clients

- **Geth**: Go Ethereum client
- **Reth**: Rust Ethereum client

### Configuration Options

```yaml
execution_client:
  name: geth  # or reth
  port: [8545, 8546]  # [HTTP, WebSocket]
  execution_type: full  # full, archive, light
  additional_args: []  # Custom command line arguments
```

### Execution Types

| Type | Description | Storage | Sync Speed |
|------|-------------|---------|------------|
| `full` | Full node with recent state | ~1.5 TB | Fast |
| `archive` | Complete historical data | ~3+ TB | Slow |
| `light` | Minimal data, relies on peers | ~50 GB | Very Fast |

### Port Configuration

```bash
# Set custom ports
starknode-kit config el client=geth port=8545,8546

# Default ports:
# Geth: 8545 (HTTP), 8546 (WebSocket)
# Reth: 8545 (HTTP), 8546 (WebSocket)
```

## Consensus Client Configuration

### Supported Clients

- **Lighthouse**: Rust-based consensus client
- **Prysm**: Go-based consensus client

### Configuration Options

```yaml
consensus_client:
  name: lighthouse  # or prysm
  port: [9000, 9001]  # [P2P, HTTP API]
  consensus_checkpoint: ""  # Optional checkpoint URL
  additional_args: []  # Custom command line arguments
```

### Checkpoint Sync

Speed up initial sync using a trusted checkpoint:

```yaml
consensus_client:
  name: lighthouse
  consensus_checkpoint: "https://mainnet.checkpoint.sigp.io"
```

**Available checkpoints**:
- Lighthouse: `https://mainnet.checkpoint.sigp.io`
- Prysm: `https://beaconstate.info/`

### Port Configuration

```bash
# Set custom ports
starknode-kit config cl client=lighthouse port=9000,9001

# Default ports:
# Lighthouse: 9000 (P2P), 9001 (HTTP API)
# Prysm: 13000 (P2P), 3500 (HTTP API)
```

## Starknet Client Configuration

### Supported Clients

- **Juno**: Go-based Starknet full node

### Configuration Options

```yaml
starknet_client:
  name: juno
  port: [6060, 6061]  # [HTTP, WebSocket]
  eth_node_url: "http://localhost:8545"  # Ethereum node URL
  additional_args: []  # Custom command line arguments
```

### Ethereum Node Connection

Juno requires an Ethereum node to verify L1 state:

```yaml
starknet_client:
  name: juno
  eth_node_url: "http://localhost:8545"  # Local Geth/Reth
  # or
  eth_node_url: "https://eth-mainnet.alchemyapi.io/v2/YOUR_KEY"  # External provider
```

### Port Configuration

```bash
# Set custom ports
starknode-kit config starknet client=juno port=6060,6061

# Default ports:
# Juno: 6060 (HTTP), 6061 (WebSocket)
```

## Advanced Configuration

### Custom Command Line Arguments

Add custom arguments to any client:

```yaml
execution_client:
  name: geth
  additional_args:
    - "--maxpeers"
    - "50"
    - "--cache"
    - "4096"

consensus_client:
  name: lighthouse
  additional_args:
    - "--max-peers"
    - "50"
    - "--target-peers"
    - "25"
```

### Environment Variables

Set environment variables for clients:

```yaml
execution_client:
  name: geth
  env_vars:
    GETH_CACHE: "4096"
    GETH_MAXPEERS: "50"
```

### Resource Limits

Configure resource limits:

```yaml
execution_client:
  name: geth
  resource_limits:
    memory: "8Gi"
    cpu: "4"
```

## Configuration Management

### View Current Configuration

```bash
# Show current config
starknode-kit config

# Show specific client config
starknode-kit config el
starknode-kit config cl
starknode-kit config starknet
```

### Edit Configuration

```bash
# Edit config file directly
nano ~/.starknode-kit/config.yaml

# Or use command line
starknode-kit config -n sepolia
starknode-kit config el client=reth port=8545,8546
```

### Validate Configuration

```bash
# Check if config is valid
starknode-kit config --validate
```

## Configuration Examples

### Development Setup

```yaml
network: sepolia
execution_client:
  name: geth
  port: [8545, 8546]
  execution_type: full
consensus_client:
  name: lighthouse
  port: [9000, 9001]
starknet_client:
  name: juno
  port: [6060, 6061]
  eth_node_url: "http://localhost:8545"
```

### Production Setup

```yaml
network: mainnet
execution_client:
  name: reth
  port: [8545, 8546]
  execution_type: full
  additional_args:
    - "--max-peers"
    - "100"
consensus_client:
  name: lighthouse
  port: [9000, 9001]
  consensus_checkpoint: "https://mainnet.checkpoint.sigp.io"
  additional_args:
    - "--max-peers"
    - "50"
starknet_client:
  name: juno
  port: [6060, 6061]
  eth_node_url: "http://localhost:8545"
```

### Archive Node Setup

```yaml
network: mainnet
execution_client:
  name: geth
  port: [8545, 8546]
  execution_type: archive
  additional_args:
    - "--cache"
    - "8192"
consensus_client:
  name: lighthouse
  port: [9000, 9001]
  execution_type: archive
```

## Troubleshooting Configuration

### Common Issues

**Invalid YAML Syntax**
```bash
# Validate YAML syntax
starknode-kit config --validate
```

**Port Conflicts**
```bash
# Check port usage
netstat -tulpn | grep :8545
```

**Missing Dependencies**
```bash
# Check if clients are installed
starknode-kit config
```

### Configuration Validation

The configuration is validated when you run commands. Common validation errors:

- Invalid network name
- Unsupported client combinations
- Port conflicts
- Missing required fields

## Best Practices

### Security
- Use non-default ports in production
- Restrict RPC access to localhost
- Use strong authentication for external access

### Performance
- Match execution type to your use case
- Use checkpoint sync for faster initial sync
- Configure appropriate resource limits

### Reliability
- Use stable network configurations
- Set up monitoring for configuration changes
- Keep backups of working configurations

## Getting Help

For configuration issues:

- Check the [Troubleshooting Guide](../operations/troubleshooting.md)
- Review client-specific documentation
- Join our [Telegram community](https://t.me/+SCPbza9fk8dkYWI0)
- Open an issue on [GitHub](https://github.com/thebuidl-grid/starknode-kit/issues)
