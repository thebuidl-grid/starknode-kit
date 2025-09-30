# Configuration Schema

Complete reference for the Starknode Kit configuration file format and all available options.

## Configuration File Location

- **Linux/macOS**: `~/.starknode-kit/config.yaml`
- **Windows**: `%USERPROFILE%\.starknode-kit\config.yaml`

## Schema Overview

```yaml
# Network configuration
network: string

# Execution client configuration
execution_client:
  name: string
  port: [number, number]
  execution_type: string
  additional_args: [string]
  env_vars: {string: string}
  resource_limits:
    memory: string
    cpu: string

# Consensus client configuration
consensus_client:
  name: string
  port: [number, number]
  consensus_checkpoint: string
  additional_args: [string]
  env_vars: {string: string}
  resource_limits:
    memory: string
    cpu: string

# Starknet client configuration
starknet_client:
  name: string
  port: [number, number]
  eth_node_url: string
  additional_args: [string]
  env_vars: {string: string}
  resource_limits:
    memory: string
    cpu: string
```

## Network Configuration

### `network`

**Type**: `string`  
**Required**: Yes  
**Default**: `mainnet`

**Valid Values**:
- `mainnet` - Ethereum mainnet
- `sepolia` - Sepolia testnet
- `holesky` - Holesky testnet
- `goerli` - Goerli testnet (deprecated)

**Example**:
```yaml
network: mainnet
```

## Execution Client Configuration

### `execution_client`

**Type**: `object`  
**Required**: No

#### `execution_client.name`

**Type**: `string`  
**Required**: Yes

**Valid Values**:
- `geth` - Go Ethereum client
- `reth` - Rust Ethereum client

**Example**:
```yaml
execution_client:
  name: geth
```

#### `execution_client.port`

**Type**: `[number, number]`  
**Required**: No  
**Default**: `[8545, 8546]`

**Description**: HTTP and WebSocket RPC ports

**Example**:
```yaml
execution_client:
  port: [8545, 8546]
```

#### `execution_client.execution_type`

**Type**: `string`  
**Required**: No  
**Default**: `full`

**Valid Values**:
- `full` - Full node with recent state
- `archive` - Complete historical data
- `light` - Minimal data, relies on peers

**Example**:
```yaml
execution_client:
  execution_type: full
```

#### `execution_client.additional_args`

**Type**: `[string]`  
**Required**: No  
**Default**: `[]`

**Description**: Additional command line arguments

**Example**:
```yaml
execution_client:
  additional_args:
    - "--maxpeers"
    - "50"
    - "--cache"
    - "4096"
```

#### `execution_client.env_vars`

**Type**: `{string: string}`  
**Required**: No  
**Default**: `{}`

**Description**: Environment variables for the client

**Example**:
```yaml
execution_client:
  env_vars:
    GETH_CACHE: "4096"
    GETH_MAXPEERS: "50"
```

#### `execution_client.resource_limits`

**Type**: `object`  
**Required**: No

##### `execution_client.resource_limits.memory`

**Type**: `string`  
**Required**: No

**Description**: Memory limit (e.g., "8Gi", "4G")

**Example**:
```yaml
execution_client:
  resource_limits:
    memory: "8Gi"
```

##### `execution_client.resource_limits.cpu`

**Type**: `string`  
**Required**: No

**Description**: CPU limit (e.g., "4", "2.5")

**Example**:
```yaml
execution_client:
  resource_limits:
    cpu: "4"
```

## Consensus Client Configuration

### `consensus_client`

**Type**: `object`  
**Required**: No

#### `consensus_client.name`

**Type**: `string`  
**Required**: Yes

**Valid Values**:
- `lighthouse` - Rust-based consensus client
- `prysm` - Go-based consensus client

**Example**:
```yaml
consensus_client:
  name: lighthouse
```

#### `consensus_client.port`

**Type**: `[number, number]`  
**Required**: No  
**Default**: `[9000, 9001]` (Lighthouse), `[13000, 3500]` (Prysm)

**Description**: P2P and HTTP API ports

**Example**:
```yaml
consensus_client:
  port: [9000, 9001]
```

#### `consensus_client.consensus_checkpoint`

**Type**: `string`  
**Required**: No  
**Default**: `""`

**Description**: Checkpoint URL for faster sync

**Example**:
```yaml
consensus_client:
  consensus_checkpoint: "https://mainnet.checkpoint.sigp.io"
```

#### `consensus_client.additional_args`

**Type**: `[string]`  
**Required**: No  
**Default**: `[]`

**Description**: Additional command line arguments

**Example**:
```yaml
consensus_client:
  additional_args:
    - "--max-peers"
    - "50"
    - "--target-peers"
    - "25"
```

#### `consensus_client.env_vars`

**Type**: `{string: string}`  
**Required**: No  
**Default**: `{}`

**Description**: Environment variables for the client

**Example**:
```yaml
consensus_client:
  env_vars:
    LIGHTHOUSE_MAX_PEERS: "50"
    LIGHTHOUSE_TARGET_PEERS: "25"
```

#### `consensus_client.resource_limits`

**Type**: `object`  
**Required**: No

##### `consensus_client.resource_limits.memory`

**Type**: `string`  
**Required**: No

**Description**: Memory limit (e.g., "4Gi", "2G")

**Example**:
```yaml
consensus_client:
  resource_limits:
    memory: "4Gi"
```

##### `consensus_client.resource_limits.cpu`

**Type**: `string`  
**Required**: No

**Description**: CPU limit (e.g., "2", "1.5")

**Example**:
```yaml
consensus_client:
  resource_limits:
    cpu: "2"
```

## Starknet Client Configuration

### `starknet_client`

**Type**: `object`  
**Required**: No

#### `starknet_client.name`

**Type**: `string`  
**Required**: Yes

**Valid Values**:
- `juno` - Go-based Starknet full node

**Example**:
```yaml
starknet_client:
  name: juno
```

#### `starknet_client.port`

**Type**: `[number, number]`  
**Required**: No  
**Default**: `[6060, 6061]`

**Description**: HTTP and WebSocket RPC ports

**Example**:
```yaml
starknet_client:
  port: [6060, 6061]
```

#### `starknet_client.eth_node_url`

**Type**: `string`  
**Required**: Yes

**Description**: Ethereum node URL for L1 verification

**Example**:
```yaml
starknet_client:
  eth_node_url: "http://localhost:8545"
```

#### `starknet_client.additional_args`

**Type**: `[string]`  
**Required**: No  
**Default**: `[]`

**Description**: Additional command line arguments

**Example**:
```yaml
starknet_client:
  additional_args:
    - "--max-peers"
    - "50"
    - "--log-level"
    - "info"
```

#### `starknet_client.env_vars`

**Type**: `{string: string}`  
**Required**: No  
**Default**: `{}`

**Description**: Environment variables for the client

**Example**:
```yaml
starknet_client:
  env_vars:
    JUNO_MAX_PEERS: "50"
    JUNO_LOG_LEVEL: "info"
```

#### `starknet_client.resource_limits`

**Type**: `object`  
**Required**: No

##### `starknet_client.resource_limits.memory`

**Type**: `string`  
**Required**: No

**Description**: Memory limit (e.g., "4Gi", "2G")

**Example**:
```yaml
starknet_client:
  resource_limits:
    memory: "4Gi"
```

##### `starknet_client.resource_limits.cpu`

**Type**: `string`  
**Required**: No

**Description**: CPU limit (e.g., "2", "1.5")

**Example**:
```yaml
starknet_client:
  resource_limits:
    cpu: "2"
```

## Complete Configuration Examples

### Development Setup

```yaml
network: sepolia
execution_client:
  name: geth
  port: [8545, 8546]
  execution_type: full
  additional_args:
    - "--maxpeers"
    - "25"
    - "--cache"
    - "2048"
consensus_client:
  name: lighthouse
  port: [9000, 9001]
  additional_args:
    - "--max-peers"
    - "25"
starknet_client:
  name: juno
  port: [6060, 6061]
  eth_node_url: "http://localhost:8545"
  additional_args:
    - "--max-peers"
    - "25"
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
    - "--cache"
    - "8192"
  resource_limits:
    memory: "16Gi"
    cpu: "8"
consensus_client:
  name: lighthouse
  port: [9000, 9001]
  consensus_checkpoint: "https://mainnet.checkpoint.sigp.io"
  additional_args:
    - "--max-peers"
    - "50"
    - "--target-peers"
    - "25"
  resource_limits:
    memory: "8Gi"
    cpu: "4"
starknet_client:
  name: juno
  port: [6060, 6061]
  eth_node_url: "http://localhost:8545"
  additional_args:
    - "--max-peers"
    - "50"
  resource_limits:
    memory: "8Gi"
    cpu: "4"
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
    - "16384"
    - "--maxpeers"
    - "100"
  resource_limits:
    memory: "32Gi"
    cpu: "16"
consensus_client:
  name: lighthouse
  port: [9000, 9001]
  execution_type: archive
  additional_args:
    - "--max-peers"
    - "50"
  resource_limits:
    memory: "16Gi"
    cpu: "8"
```

## Validation Rules

### Required Fields

- `network` - Must be a valid network name
- `execution_client.name` - Must be a valid execution client
- `consensus_client.name` - Must be a valid consensus client
- `starknet_client.name` - Must be a valid Starknet client
- `starknet_client.eth_node_url` - Must be a valid URL

### Port Validation

- Ports must be in range 1024-65535
- Ports must not conflict with system ports
- HTTP and WebSocket ports must be different

### Resource Limits

- Memory must be in format "XGi" or "XG"
- CPU must be a positive number
- Limits must not exceed system resources

## Environment Variables

### Configuration Override

Environment variables can override configuration file settings:

```bash
# Override network
export STARKNODE_NETWORK=sepolia

# Override execution client
export STARKNODE_EXECUTION_CLIENT=geth

# Override consensus client
export STARKNODE_CONSENSUS_CLIENT=lighthouse
```

### Client-Specific Variables

```bash
# Geth variables
export GETH_CACHE=4096
export GETH_MAXPEERS=50

# Lighthouse variables
export LIGHTHOUSE_MAX_PEERS=50
export LIGHTHOUSE_TARGET_PEERS=25

# Juno variables
export JUNO_MAX_PEERS=50
export JUNO_LOG_LEVEL=info
```

## Configuration Management

### Validation

```bash
# Validate configuration
starknode-kit config --validate
```

### Backup

```bash
# Backup configuration
cp ~/.starknode-kit/config.yaml ~/.starknode-kit/config.yaml.backup
```

### Reset

```bash
# Reset to defaults
rm ~/.starknode-kit/config.yaml
starknode-kit init
```

## Getting Help

For configuration issues:

- Check the [Troubleshooting Guide](../operations/troubleshooting.md)
- Review client-specific documentation
- Join our [Telegram community](https://t.me/+SCPbza9fk8dkYWI0)
- Open an issue on [GitHub](https://github.com/thebuidl-grid/starknode-kit/issues)
