# Command Reference

Complete reference for all Starknode Kit commands, flags, and options.

## Global Commands

### `starknode-kit --help`

Display help information for the Starknode Kit CLI.

```bash
starknode-kit --help
```

**Output**: Shows available commands and global flags.

---

## Core Commands

### `init`

Create a default configuration file.

```bash
starknode-kit init
```

**Description**: Initializes a new configuration file with default settings.

**Files Created**:
- `~/.starknode-kit/config.yaml`

**Example**:
```bash
starknode-kit init
# Creates default configuration file
```

---

### `add`

Add an Ethereum or Starknet client to the configuration.

```bash
starknode-kit add [flags]
```

**Flags**:
- `-e, --execution_client string`: Execution client (geth, reth)
- `-c, --consensus_client string`: Consensus client (lighthouse, prysm)
- `-s, --starknet_client string`: Starknet client (juno)

**Examples**:
```bash
# Add execution client
starknode-kit add --execution_client geth

# Add consensus client
starknode-kit add --consensus_client lighthouse

# Add Starknet client
starknode-kit add --starknet_client juno

# Add client pair
starknode-kit add --execution_client geth --consensus_client lighthouse

# Add all clients
starknode-kit add --execution_client geth --consensus_client lighthouse --starknet_client juno
```

---

### `config`

Show and modify the current configuration.

```bash
starknode-kit config [flags] [options]
```

**Flags**:
- `-n, --network string`: Set network (mainnet, sepolia, holesky)
- `--validate`: Validate configuration file

**Options**:
- `el client=<name> port=<ports>`: Configure execution client
- `cl client=<name> port=<ports>`: Configure consensus client
- `starknet client=<name> port=<ports>`: Configure Starknet client

**Examples**:
```bash
# Show current configuration
starknode-kit config

# Set network
starknode-kit config -n sepolia

# Configure execution client
starknode-kit config el client=geth port=8545,8546

# Configure consensus client
starknode-kit config cl client=lighthouse port=9000,9001

# Validate configuration
starknode-kit config --validate
```

---

### `start`

Start the configured Ethereum clients.

```bash
starknode-kit start
```

**Description**: Starts execution and consensus clients. Does not start Starknet clients.

**Examples**:
```bash
# Start all configured clients
starknode-kit start
```

**Note**: Use `starknode-kit run juno` to start Starknet clients.

---

### `stop`

Stop all running clients.

```bash
starknode-kit stop
```

**Description**: Gracefully stops all running clients.

**Examples**:
```bash
# Stop all clients
starknode-kit stop
```

---

### `run`

Run local Starknet infrastructure services.

```bash
starknode-kit run <service>
```

**Subcommands**:
- `juno`: Run a Juno Starknet node

**Examples**:
```bash
# Run Juno node
starknode-kit run juno
```

---

### `monitor`

Launch the real-time monitoring dashboard.

```bash
starknode-kit monitor [flags]
```

**Flags**:
- `--interval duration`: Refresh interval (default: 5s)
- `--basic`: Basic monitoring mode
- `--restart`: Restart monitoring service

**Examples**:
```bash
# Start monitoring dashboard
starknode-kit monitor

# Custom refresh interval
starknode-kit monitor --interval 10s

# Basic monitoring mode
starknode-kit monitor --basic
```

---

### `remove`

Remove a configured client.

```bash
starknode-kit remove [flags]
```

**Flags**:
- `-e, --execution_client string`: Remove execution client
- `-c, --consensus_client string`: Remove consensus client
- `-s, --starknet_client string`: Remove Starknet client

**Examples**:
```bash
# Remove execution client
starknode-kit remove --execution_client geth

# Remove consensus client
starknode-kit remove --consensus_client lighthouse

# Remove Starknet client
starknode-kit remove --starknet_client juno
```

---

### `update`

Check for and install client updates.

```bash
starknode-kit update [flags]
```

**Flags**:
- `-e, --execution_client string`: Update specific execution client
- `-c, --consensus_client string`: Update specific consensus client
- `-s, --starknet_client string`: Update specific Starknet client

**Examples**:
```bash
# Update all clients
starknode-kit update

# Update specific client
starknode-kit update --execution_client geth
starknode-kit update --consensus_client lighthouse
starknode-kit update --starknet_client juno
```

---

## Utility Commands

### `completion`

Generate shell completion scripts.

```bash
starknode-kit completion <shell>
```

**Supported Shells**:
- `bash`
- `zsh`
- `fish`
- `powershell`

**Examples**:
```bash
# Generate bash completion
starknode-kit completion bash

# Install bash completion
starknode-kit completion bash > /etc/bash_completion.d/starknode-kit

# Generate zsh completion
starknode-kit completion zsh > ~/.zsh/completions/_starknode-kit
```

---

## Global Flags

### `-h, --help`

Show help information for any command.

```bash
starknode-kit <command> --help
```

**Examples**:
```bash
starknode-kit add --help
starknode-kit config --help
starknode-kit start --help
```

---

## Command Combinations

### Complete Setup Workflow

```bash
# 1. Initialize configuration
starknode-kit init

# 2. Add clients
starknode-kit add --execution_client geth --consensus_client lighthouse --starknet_client juno

# 3. Configure network
starknode-kit config -n mainnet

# 4. Start Ethereum clients
starknode-kit start

# 5. Start Starknet client (in separate terminal)
starknode-kit run juno

# 6. Monitor everything
starknode-kit monitor
```

### Development Workflow

```bash
# 1. Set up testnet
starknode-kit config -n sepolia

# 2. Add development clients
starknode-kit add --execution_client geth --consensus_client prysm

# 3. Start clients
starknode-kit start

# 4. Monitor progress
starknode-kit monitor
```

### Maintenance Workflow

```bash
# 1. Check for updates
starknode-kit update

# 2. Stop clients
starknode-kit stop

# 3. Update clients
starknode-kit update

# 4. Restart clients
starknode-kit start

# 5. Verify status
starknode-kit monitor
```

---

## Exit Codes

| Code | Description |
|------|-------------|
| `0` | Success |
| `1` | General error |
| `2` | Configuration error |
| `3` | Client error |
| `4` | Network error |
| `5` | Permission error |

---

## Environment Variables

### Configuration

- `STARKNODE_CONFIG_PATH`: Custom config file path
- `STARKNODE_LOG_LEVEL`: Log level (debug, info, warn, error)
- `STARKNODE_DATA_DIR`: Custom data directory

### Client Configuration

- `GETH_CACHE`: Geth cache size
- `LIGHTHOUSE_MAX_PEERS`: Lighthouse max peers
- `JUNO_ETH_NODE`: Juno Ethereum node URL

**Examples**:
```bash
# Set custom config path
export STARKNODE_CONFIG_PATH=/custom/path/config.yaml

# Set log level
export STARKNODE_LOG_LEVEL=debug

# Set Geth cache
export GETH_CACHE=4096
```

---

## Configuration File Format

The configuration file is in YAML format:

```yaml
network: mainnet
execution_client:
  name: geth
  port: [8545, 8546]
  execution_type: full
  additional_args: []
consensus_client:
  name: lighthouse
  port: [9000, 9001]
  consensus_checkpoint: ""
  additional_args: []
starknet_client:
  name: juno
  port: [6060, 6061]
  eth_node_url: "http://localhost:8545"
  additional_args: []
```

---

## Troubleshooting Commands

### Debug Information

```bash
# Show detailed configuration
starknode-kit config --validate

# Check client status
starknode-kit config

# Monitor with debug logs
starknode-kit monitor --interval 1s
```

### Log Analysis

```bash
# View client logs
tail -f ~/.starknode-kit/logs/geth.log
tail -f ~/.starknode-kit/logs/lighthouse.log
tail -f ~/.starknode-kit/logs/juno.log

# Search for errors
grep -i error ~/.starknode-kit/logs/*.log
```

---

## Getting Help

For command-specific help:

```bash
# General help
starknode-kit --help

# Command help
starknode-kit <command> --help

# Subcommand help
starknode-kit run --help
```

For additional support:

- Check the [Troubleshooting Guide](../operations/troubleshooting.md)
- Join our [Telegram community](https://t.me/+SCPbza9fk8dkYWI0)
- Open an issue on [GitHub](https://github.com/thebuidl-grid/starknode-kit/issues)
