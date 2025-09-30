# Running Nodes

This guide covers how to run and manage your Ethereum and Starknet nodes using Starknode Kit.

## Starting Your Nodes

### Start Ethereum Clients

Start your configured execution and consensus clients:

```bash
starknode-kit start
```

This command will:
1. Check your configuration
2. Start the execution client first
3. Start the consensus client
4. Begin the synchronization process

### Start Starknet Client

Start your Juno Starknet node:

```bash
starknode-kit run juno
```

**Note**: Juno requires an Ethereum node connection to verify L1 state.

## Node Startup Process

### Execution Client Startup

**Geth**:
```bash
# Geth starts with these default settings:
# - HTTP RPC on port 8545
# - WebSocket RPC on port 8546
# - P2P networking on port 30303
# - Full sync mode
```

**Reth**:
```bash
# Reth starts with these default settings:
# - HTTP RPC on port 8545
# - WebSocket RPC on port 8546
# - P2P networking on port 30303
# - Full sync mode
```

### Consensus Client Startup

**Lighthouse**:
```bash
# Lighthouse starts with these default settings:
# - P2P networking on port 9000
# - HTTP API on port 9001
# - Beacon node mode
# - Mainnet network
```

**Prysm**:
```bash
# Prysm starts with these default settings:
# - P2P networking on port 13000
# - HTTP API on port 3500
# - Beacon node mode
# - Mainnet network
```

### Starknet Client Startup

**Juno**:
```bash
# Juno starts with these default settings:
# - HTTP RPC on port 6060
# - WebSocket RPC on port 6061
# - Ethereum node connection required
# - Mainnet network
```

## Synchronization Process

### Initial Sync

When you first start your nodes, they will begin synchronizing with the network:

**Execution Client Sync**:
- Downloads and verifies all blocks
- Builds the state trie
- Can take 1-3 days depending on hardware

**Consensus Client Sync**:
- Downloads beacon chain data
- Verifies consensus rules
- Typically faster than execution sync

**Starknet Client Sync**:
- Syncs L2 state from L1
- Requires Ethereum node connection
- Usually completes in 2-6 hours

### Sync Progress Monitoring

Monitor sync progress using the built-in dashboard:

```bash
starknode-kit monitor
```

Key metrics to watch:
- **Block height**: Current synced block number
- **Sync percentage**: Progress towards full sync
- **Peer count**: Number of connected peers
- **Sync speed**: Blocks per second

### Checkpoint Sync

Speed up initial sync using trusted checkpoints:

```bash
# Configure checkpoint sync for Lighthouse
starknode-kit config cl client=lighthouse consensus_checkpoint="https://mainnet.checkpoint.sigp.io"

# Configure checkpoint sync for Prysm
starknode-kit config cl client=prysm consensus_checkpoint="https://beaconstate.info/"
```

## Node Management

### Stop Nodes

Stop all running clients:

```bash
starknode-kit stop
```

This gracefully shuts down all clients.

### Restart Nodes

Restart all clients:

```bash
starknode-kit stop
starknode-kit start
```

### Check Node Status

Check if nodes are running:

```bash
# Check configuration
starknode-kit config

# Check running processes
ps aux | grep -E "(geth|lighthouse|juno)"

# Check ports
netstat -tulpn | grep -E "(8545|9000|6060)"
```

## Network Configuration

### Mainnet

Run on Ethereum mainnet:

```bash
starknode-kit config -n mainnet
starknode-kit start
```

### Testnets

Run on testnets for development:

```bash
# Sepolia testnet
starknode-kit config -n sepolia
starknode-kit start

# Holesky testnet
starknode-kit config -n holesky
starknode-kit start
```

### Custom Networks

Configure custom network settings:

```yaml
# config.yaml
network: custom
execution_client:
  name: geth
  additional_args:
    - "--networkid"
    - "12345"
    - "--genesis"
    - "/path/to/genesis.json"
```

## Performance Optimization

### Resource Allocation

Optimize resource usage for your hardware:

```bash
# High-memory system (32GB+)
starknode-kit config el client=geth additional_args="--cache,8192,--maxpeers,100"

# Low-memory system (16GB)
starknode-kit config el client=geth additional_args="--cache,2048,--maxpeers,25"

# High-performance system
starknode-kit config el client=reth additional_args="--max-peers,100"
```

### Network Optimization

Optimize network settings:

```bash
# Increase peer connections
starknode-kit config el client=geth additional_args="--maxpeers,100"

# Optimize for low bandwidth
starknode-kit config el client=geth additional_args="--maxpeers,25,--lightserv,0"
```

### Storage Optimization

Optimize storage usage:

```bash
# Use light client mode
starknode-kit config el client=geth execution_type=light

# Use archive mode for complete data
starknode-kit config el client=geth execution_type=archive
```

## Monitoring and Logs

### Real-time Monitoring

Use the built-in monitoring dashboard:

```bash
starknode-kit monitor
```

Features:
- Real-time sync progress
- Resource usage statistics
- Peer connection status
- Live log streaming

### Log Files

Access client log files:

```bash
# Geth logs
tail -f ~/.starknode-kit/logs/geth.log

# Lighthouse logs
tail -f ~/.starknode-kit/logs/lighthouse.log

# Juno logs
tail -f ~/.starknode-kit/logs/juno.log
```

### Log Analysis

Analyze logs for issues:

```bash
# Search for errors
grep -i error ~/.starknode-kit/logs/*.log

# Monitor sync progress
grep -i "sync" ~/.starknode-kit/logs/geth.log

# Check peer connections
grep -i "peer" ~/.starknode-kit/logs/*.log
```

## Troubleshooting

### Common Issues

**Nodes Won't Start**
```bash
# Check configuration
starknode-kit config --validate

# Check port conflicts
netstat -tulpn | grep -E "(8545|9000|6060)"

# Check disk space
df -h
```

**Sync Stalls**
```bash
# Restart nodes
starknode-kit stop && starknode-kit start

# Check network connectivity
ping 8.8.8.8

# Use checkpoint sync
starknode-kit config cl client=lighthouse consensus_checkpoint="https://mainnet.checkpoint.sigp.io"
```

**High Resource Usage**
```bash
# Check resource usage
htop
iotop

# Reduce resource allocation
starknode-kit config el client=geth additional_args="--cache,2048,--maxpeers,25"
```

### Performance Issues

**Slow Sync**
- Check disk I/O performance
- Verify network connection
- Consider using checkpoint sync
- Optimize client settings

**High Memory Usage**
- Reduce cache size
- Use light client mode
- Add swap space
- Monitor for memory leaks

**High CPU Usage**
- Reduce peer count
- Use lighter execution type
- Check for stuck processes
- Monitor system load

## Best Practices

### Security

- Run nodes behind a firewall
- Use non-default ports in production
- Restrict RPC access to localhost
- Keep clients updated

### Reliability

- Use stable network configurations
- Monitor system resources
- Set up automated restarts
- Keep backups of configuration

### Performance

- Match execution type to use case
- Optimize for your hardware
- Use checkpoint sync for faster initial sync
- Monitor and tune settings

## Getting Help

For node operation issues:

- Check the [Troubleshooting Guide](troubleshooting.md)
- Review client-specific documentation
- Join our [Telegram community](https://t.me/+SCPbza9fk8dkYWI0)
- Open an issue on [GitHub](https://github.com/thebuidl-grid/starknode-kit/issues)
