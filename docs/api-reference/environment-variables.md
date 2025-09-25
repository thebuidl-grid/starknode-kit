# Environment Variables

Complete reference for all environment variables used by Starknode Kit and its clients.

## Starknode Kit Variables

### Configuration

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `STARKNODE_CONFIG_PATH` | Custom config file path | `~/.starknode-kit/config.yaml` | `/custom/path/config.yaml` |
| `STARKNODE_LOG_LEVEL` | Log level for Starknode Kit | `info` | `debug`, `info`, `warn`, `error` |
| `STARKNODE_DATA_DIR` | Custom data directory | `~/.starknode-kit/` | `/custom/data/path/` |
| `STARKNODE_NETWORK` | Override network setting | From config | `mainnet`, `sepolia`, `holesky` |
| `STARKNODE_EXECUTION_CLIENT` | Override execution client | From config | `geth`, `reth` |
| `STARKNODE_CONSENSUS_CLIENT` | Override consensus client | From config | `lighthouse`, `prysm` |
| `STARKNODE_STARKNET_CLIENT` | Override Starknet client | From config | `juno` |

### Monitoring

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `STARKNODE_MONITOR_INTERVAL` | Monitoring refresh interval | `5s` | `10s`, `30s` |
| `STARKNODE_MONITOR_LOG_LEVEL` | Monitoring log level | `info` | `debug`, `info`, `warn` |
| `STARKNODE_MONITOR_METRICS` | Enable metrics collection | `true` | `true`, `false` |
| `STARKNODE_MONITOR_BASIC` | Use basic monitoring mode | `false` | `true`, `false` |

## Execution Client Variables

### Geth Variables

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `GETH_CACHE` | Memory cache size (MB) | `1024` | `4096`, `8192` |
| `GETH_MAXPEERS` | Maximum number of peers | `50` | `100`, `25` |
| `GETH_LIGHTKDF` | Use light KDF for key derivation | `false` | `true`, `false` |
| `GETH_NETWORKID` | Network ID | `1` (mainnet) | `11155111` (sepolia) |
| `GETH_DATADIR` | Data directory | `~/.starknode-kit/data/geth/` | `/custom/path/` |
| `GETH_HTTP_PORT` | HTTP RPC port | `8545` | `8547` |
| `GETH_WS_PORT` | WebSocket RPC port | `8546` | `8548` |
| `GETH_P2P_PORT` | P2P networking port | `30303` | `30305` |
| `GETH_HTTP_ADDR` | HTTP RPC address | `localhost` | `0.0.0.0` |
| `GETH_WS_ADDR` | WebSocket RPC address | `localhost` | `0.0.0.0` |
| `GETH_HTTP_API` | HTTP RPC API modules | `eth,net,web3` | `eth,net,web3,admin` |
| `GETH_WS_API` | WebSocket RPC API modules | `eth,net,web3` | `eth,net,web3,admin` |
| `GETH_HTTP_CORS` | HTTP RPC CORS origins | `*` | `http://localhost:3000` |
| `GETH_WS_ORIGINS` | WebSocket RPC origins | `*` | `http://localhost:3000` |
| `GETH_HTTP_VHOSTS` | HTTP RPC virtual hosts | `localhost` | `localhost,example.com` |
| `GETH_WS_VHOSTS` | WebSocket RPC virtual hosts | `localhost` | `localhost,example.com` |
| `GETH_METRICS` | Enable metrics | `false` | `true`, `false` |
| `GETH_METRICS_ADDR` | Metrics address | `127.0.0.1` | `0.0.0.0` |
| `GETH_METRICS_PORT` | Metrics port | `6060` | `6061` |
| `GETH_PPROF` | Enable pprof | `false` | `true`, `false` |
| `GETH_PPROF_ADDR` | pprof address | `127.0.0.1` | `0.0.0.0` |
| `GETH_PPROF_PORT` | pprof port | `6060` | `6061` |

### Reth Variables

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `RETH_DATADIR` | Data directory | `~/.starknode-kit/data/reth/` | `/custom/path/` |
| `RETH_HTTP_PORT` | HTTP RPC port | `8545` | `8547` |
| `RETH_WS_PORT` | WebSocket RPC port | `8546` | `8548` |
| `RETH_P2P_PORT` | P2P networking port | `30303` | `30305` |
| `RETH_HTTP_ADDR` | HTTP RPC address | `127.0.0.1` | `0.0.0.0` |
| `RETH_WS_ADDR` | WebSocket RPC address | `127.0.0.1` | `0.0.0.0` |
| `RETH_MAX_PEERS` | Maximum number of peers | `50` | `100`, `25` |
| `RETH_TARGET_PEERS` | Target number of peers | `25` | `50`, `15` |
| `RETH_CACHE_SIZE` | Cache size (MB) | `1024` | `4096`, `8192` |
| `RETH_NETWORK` | Network name | `mainnet` | `sepolia`, `holesky` |
| `RETH_CHAIN` | Chain specification | `mainnet` | `sepolia`, `holesky` |
| `RETH_HTTP_CORS` | HTTP RPC CORS origins | `*` | `http://localhost:3000` |
| `RETH_WS_ORIGINS` | WebSocket RPC origins | `*` | `http://localhost:3000` |
| `RETH_METRICS` | Enable metrics | `false` | `true`, `false` |
| `RETH_METRICS_ADDR` | Metrics address | `127.0.0.1` | `0.0.0.0` |
| `RETH_METRICS_PORT` | Metrics port | `9001` | `9002` |

## Consensus Client Variables

### Lighthouse Variables

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `LIGHTHOUSE_DATADIR` | Data directory | `~/.starknode-kit/data/lighthouse/` | `/custom/path/` |
| `LIGHTHOUSE_P2P_PORT` | P2P networking port | `9000` | `9002` |
| `LIGHTHOUSE_HTTP_PORT` | HTTP API port | `9001` | `9003` |
| `LIGHTHOUSE_P2P_ADDR` | P2P address | `0.0.0.0` | `127.0.0.1` |
| `LIGHTHOUSE_HTTP_ADDR` | HTTP API address | `127.0.0.1` | `0.0.0.0` |
| `LIGHTHOUSE_MAX_PEERS` | Maximum number of peers | `50` | `100`, `25` |
| `LIGHTHOUSE_TARGET_PEERS` | Target number of peers | `25` | `50`, `15` |
| `LIGHTHOUSE_NETWORK` | Network name | `mainnet` | `sepolia`, `holesky` |
| `LIGHTHOUSE_CHECKPOINT_SYNC` | Checkpoint sync URL | `""` | `https://mainnet.checkpoint.sigp.io` |
| `LIGHTHOUSE_HTTP_CORS` | HTTP API CORS origins | `*` | `http://localhost:3000` |
| `LIGHTHOUSE_METRICS` | Enable metrics | `false` | `true`, `false` |
| `LIGHTHOUSE_METRICS_ADDR` | Metrics address | `127.0.0.1` | `0.0.0.0` |
| `LIGHTHOUSE_METRICS_PORT` | Metrics port | `5054` | `5055` |
| `LIGHTHOUSE_LOG_LEVEL` | Log level | `info` | `debug`, `warn`, `error` |
| `LIGHTHOUSE_LOG_COLOR` | Colored log output | `true` | `true`, `false` |
| `LIGHTHOUSE_LOG_FILE` | Log file path | `~/.starknode-kit/logs/lighthouse.log` | `/custom/path/log.log` |

### Prysm Variables

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `PRYSM_DATADIR` | Data directory | `~/.starknode-kit/data/prysm/` | `/custom/path/` |
| `PRYSM_P2P_PORT` | P2P networking port | `13000` | `13002` |
| `PRYSM_HTTP_PORT` | HTTP API port | `3500` | `3502` |
| `PRYSM_P2P_ADDR` | P2P address | `0.0.0.0` | `127.0.0.1` |
| `PRYSM_HTTP_ADDR` | HTTP API address | `127.0.0.1` | `0.0.0.0` |
| `PRYSM_MAX_PEERS` | Maximum number of peers | `50` | `100`, `25` |
| `PRYSM_TARGET_PEERS` | Target number of peers | `25` | `50`, `15` |
| `PRYSM_NETWORK` | Network name | `mainnet` | `sepolia`, `holesky` |
| `PRYSM_CHECKPOINT_SYNC` | Checkpoint sync URL | `""` | `https://beaconstate.info/` |
| `PRYSM_HTTP_CORS` | HTTP API CORS origins | `*` | `http://localhost:3000` |
| `PRYSM_METRICS` | Enable metrics | `false` | `true`, `false` |
| `PRYSM_METRICS_ADDR` | Metrics address | `127.0.0.1` | `0.0.0.0` |
| `PRYSM_METRICS_PORT` | Metrics port | `8080` | `8081` |
| `PRYSM_LOG_LEVEL` | Log level | `info` | `debug`, `warn`, `error` |
| `PRYSM_LOG_FORMAT` | Log format | `text` | `text`, `json` |
| `PRYSM_LOG_FILE` | Log file path | `~/.starknode-kit/logs/prysm.log` | `/custom/path/log.log` |

## Starknet Client Variables

### Juno Variables

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `JUNO_DATADIR` | Data directory | `~/.starknode-kit/data/juno/` | `/custom/path/` |
| `JUNO_HTTP_PORT` | HTTP RPC port | `6060` | `6062` |
| `JUNO_WS_PORT` | WebSocket RPC port | `6061` | `6063` |
| `JUNO_HTTP_ADDR` | HTTP RPC address | `127.0.0.1` | `0.0.0.0` |
| `JUNO_WS_ADDR` | WebSocket RPC address | `127.0.0.1` | `0.0.0.0` |
| `JUNO_ETH_NODE` | Ethereum node URL | `http://localhost:8545` | `https://eth-mainnet.alchemyapi.io/v2/KEY` |
| `JUNO_MAX_PEERS` | Maximum number of peers | `50` | `100`, `25` |
| `JUNO_TARGET_PEERS` | Target number of peers | `25` | `50`, `15` |
| `JUNO_NETWORK` | Network name | `mainnet` | `sepolia`, `holesky` |
| `JUNO_CHAIN` | Chain specification | `mainnet` | `sepolia`, `holesky` |
| `JUNO_HTTP_CORS` | HTTP RPC CORS origins | `*` | `http://localhost:3000` |
| `JUNO_WS_ORIGINS` | WebSocket RPC origins | `*` | `http://localhost:3000` |
| `JUNO_METRICS` | Enable metrics | `false` | `true`, `false` |
| `JUNO_METRICS_ADDR` | Metrics address | `127.0.0.1` | `0.0.0.0` |
| `JUNO_METRICS_PORT` | Metrics port | `9090` | `9091` |
| `JUNO_LOG_LEVEL` | Log level | `info` | `debug`, `warn`, `error` |
| `JUNO_LOG_COLOR` | Colored log output | `true` | `true`, `false` |
| `JUNO_LOG_FILE` | Log file path | `~/.starknode-kit/logs/juno.log` | `/custom/path/log.log` |
| `JUNO_DB_PATH` | Database path | `~/.starknode-kit/data/juno/db/` | `/custom/db/path/` |
| `JUNO_CACHE_SIZE` | Cache size (MB) | `1024` | `4096`, `8192` |
| `JUNO_SYNC_MODE` | Sync mode | `full` | `full`, `light` |
| `JUNO_VERIFY_L1` | Verify L1 state | `true` | `true`, `false` |

## System Variables

### Resource Limits

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `STARKNODE_MEMORY_LIMIT` | Total memory limit | `16Gi` | `32Gi`, `8Gi` |
| `STARKNODE_CPU_LIMIT` | Total CPU limit | `4` | `8`, `2` |
| `STARKNODE_DISK_LIMIT` | Disk usage limit | `2Ti` | `4Ti`, `1Ti` |

### Network

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `STARKNODE_NETWORK_INTERFACE` | Network interface | `auto` | `eth0`, `wlan0` |
| `STARKNODE_BANDWIDTH_LIMIT` | Bandwidth limit | `unlimited` | `100Mbps`, `1Gbps` |
| `STARKNODE_FIREWALL` | Enable firewall | `false` | `true`, `false` |

### Security

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `STARKNODE_SSL_CERT` | SSL certificate path | `""` | `/path/to/cert.pem` |
| `STARKNODE_SSL_KEY` | SSL private key path | `""` | `/path/to/key.pem` |
| `STARKNODE_AUTH_TOKEN` | Authentication token | `""` | `your-secret-token` |
| `STARKNODE_RATE_LIMIT` | Rate limit (req/min) | `1000` | `500`, `2000` |

## Usage Examples

### Development Environment

```bash
# Set development network
export STARKNODE_NETWORK=sepolia

# Configure for development
export GETH_CACHE=2048
export GETH_MAXPEERS=25
export LIGHTHOUSE_MAX_PEERS=25
export JUNO_MAX_PEERS=25

# Enable debug logging
export STARKNODE_LOG_LEVEL=debug
export GETH_LOG_LEVEL=debug
export LIGHTHOUSE_LOG_LEVEL=debug
export JUNO_LOG_LEVEL=debug
```

### Production Environment

```bash
# Set production network
export STARKNODE_NETWORK=mainnet

# Configure for production
export GETH_CACHE=8192
export GETH_MAXPEERS=100
export LIGHTHOUSE_MAX_PEERS=50
export JUNO_MAX_PEERS=50

# Enable metrics
export GETH_METRICS=true
export LIGHTHOUSE_METRICS=true
export JUNO_METRICS=true

# Set resource limits
export STARKNODE_MEMORY_LIMIT=32Gi
export STARKNODE_CPU_LIMIT=8
```

### High-Performance Setup

```bash
# Use Reth for better performance
export STARKNODE_EXECUTION_CLIENT=reth

# Configure for high performance
export RETH_CACHE_SIZE=16384
export RETH_MAX_PEERS=100
export LIGHTHOUSE_MAX_PEERS=50
export JUNO_MAX_PEERS=50

# Enable all optimizations
export RETH_METRICS=true
export LIGHTHOUSE_METRICS=true
export JUNO_METRICS=true
```

## Environment File

Create a `.env` file in your home directory:

```bash
# ~/.starknode-kit/.env
STARKNODE_NETWORK=mainnet
STARKNODE_LOG_LEVEL=info

# Geth settings
GETH_CACHE=4096
GETH_MAXPEERS=50

# Lighthouse settings
LIGHTHOUSE_MAX_PEERS=50
LIGHTHOUSE_TARGET_PEERS=25

# Juno settings
JUNO_MAX_PEERS=50
JUNO_ETH_NODE=http://localhost:8545
```

Load environment variables:

```bash
# Load from .env file
source ~/.starknode-kit/.env

# Or use dotenv
export $(cat ~/.starknode-kit/.env | xargs)
```

## Validation

### Check Environment Variables

```bash
# List all Starknode Kit variables
env | grep STARKNODE

# List all client variables
env | grep -E "(GETH|LIGHTHOUSE|JUNO|RETH|PRYSM)"
```

### Validate Configuration

```bash
# Validate with environment variables
starknode-kit config --validate

# Check if variables are being used
starknode-kit config
```

## Best Practices

### Security

- Never commit `.env` files to version control
- Use strong authentication tokens
- Restrict network access in production
- Regularly rotate secrets

### Performance

- Set appropriate cache sizes for your hardware
- Monitor resource usage
- Adjust peer counts based on network conditions
- Use metrics to optimize settings

### Maintenance

- Document custom environment variables
- Keep variables organized in `.env` files
- Test changes in development first
- Monitor for deprecated variables

## Getting Help

For environment variable issues:

- Check the [Troubleshooting Guide](../operations/troubleshooting.md)
- Review client-specific documentation
- Join our [Telegram community](https://t.me/+SCPbza9fk8dkYWI0)
- Open an issue on [GitHub](https://github.com/thebuidl-grid/starknode-kit/issues)
