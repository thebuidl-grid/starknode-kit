# Quick Start Guide

Get your first Ethereum and Starknet node running in under 10 minutes with this step-by-step guide.

## Prerequisites

Before starting, ensure you have:
- [Starknode Kit installed](installation.md)
- [Adequate hardware](hardware-requirements.md)
- Stable internet connection

## Step 1: Initialize Configuration

Create your initial configuration file:

```bash
starknode-kit init
```

This creates a default configuration in your home directory that you can customize later.

## Step 2: Add Your First Client Pair

Add an Ethereum execution and consensus client pair:

```bash
# Add Geth (execution) + Lighthouse (consensus)
starknode-kit add --execution_client geth --consensus_client lighthouse
```

**Alternative combinations**:
```bash
# Reth + Prysm
starknode-kit add --execution_client reth --consensus_client prysm

# Geth + Prysm
starknode-kit add --execution_client geth --consensus_client prysm
```

## Step 3: Add a Starknet Client

Add Juno for Starknet support:

```bash
starknode-kit add --starknet_client juno
```

## Step 4: Verify Configuration

Check your current configuration:

```bash
starknode-kit config
```

You should see your configured clients listed.

## Step 5: Start Your Ethereum Node

Start the Ethereum execution and consensus clients:

```bash
starknode-kit start
```

This will:
- Download and install the client binaries
- Start the execution client first
- Start the consensus client
- Begin syncing with the network

## Step 6: Start Your Starknet Node

In a separate terminal, start your Juno Starknet node:

```bash
starknode-kit run juno
```

## Step 7: Monitor Your Nodes

Launch the monitoring dashboard:

```bash
starknode-kit monitor
```

This provides real-time information about:
- Client status and health
- Sync progress
- Resource usage
- Network connections

## What Happens Next

### Initial Sync Process

**Ethereum Node**:
- Execution client syncs the blockchain state
- Consensus client syncs the beacon chain
- Full sync can take 1-3 days depending on hardware

**Starknet Node (Juno)**:
- Syncs Starknet state from L1
- Requires Ethereum node connection
- Typically faster than Ethereum sync

### Expected Timeline

| Phase | Duration | Description |
|-------|----------|-------------|
| Client Installation | 5-15 minutes | Download and setup binaries |
| Ethereum Sync Start | 1-2 hours | Begin blockchain synchronization |
| Partial Sync | 6-12 hours | Sync recent blocks |
| Full Sync | 1-3 days | Complete historical sync |
| Starknet Sync | 2-6 hours | Sync Starknet state |

## Common First-Time Issues

### Client Installation Fails

**Problem**: Clients fail to download or install
**Solution**: Check internet connection and available disk space

```bash
# Check disk space
df -h

# Check internet connection
ping google.com
```

### Sync Stalls

**Problem**: Sync progress stops or becomes very slow
**Solution**: Restart the clients

```bash
# Stop clients
starknode-kit stop

# Start again
starknode-kit start
```

### Port Conflicts

**Problem**: Clients can't bind to required ports
**Solution**: Check for conflicting services

```bash
# Check port usage
netstat -tulpn | grep :8545  # Geth default port
netstat -tulpn | grep :9000  # Lighthouse default port
```

## Next Steps

Once your nodes are synced and running:

1. **Configure Monitoring**: Set up alerts and monitoring
2. **Security Hardening**: Review security best practices
3. **Performance Tuning**: Optimize for your hardware
4. **Backup Strategy**: Implement backup procedures

## Quick Commands Reference

```bash
# Check status
starknode-kit config

# Stop all clients
starknode-kit stop

# Restart clients
starknode-kit stop && starknode-kit start

# Update clients
starknode-kit update

# Remove a client
starknode-kit remove --consensus_client lighthouse

# Get help
starknode-kit --help
starknode-kit help <command>
```

## Getting Help

If you encounter issues:

1. Check the [Troubleshooting Guide](../operations/troubleshooting.md)
2. Review client-specific logs
3. Join our [Telegram community](https://t.me/+SCPbza9fk8dkYWI0)
4. Open an issue on [GitHub](https://github.com/thebuidl-grid/starknode-kit/issues)

## What's Next?

Now that you have your first node running:

- [Learn about client configuration](configuration.md)
- [Set up monitoring and alerts](../operations/monitoring.md)
- [Explore advanced features](../advanced/network-configuration.md)
- [Join the community](https://t.me/+SCPbza9fk8dkYWI0)
