# Troubleshooting

This guide helps you diagnose and resolve common issues with Starknode Kit.

## Common Issues

### Installation Problems

#### Command Not Found

**Problem**: `starknode-kit: command not found`

**Solutions**:
```bash
# Check if binary exists
ls -la /usr/local/bin/starknode-kit

# Add to PATH
echo 'export PATH=$PATH:/usr/local/bin' >> ~/.bashrc
source ~/.bashrc

# Reinstall if missing
curl -sSL https://raw.githubusercontent.com/thebuidl-grid/starknode-kit/main/install.sh | bash
```

#### Permission Denied

**Problem**: Permission denied when running commands

**Solutions**:
```bash
# Fix permissions
sudo chmod +x /usr/local/bin/starknode-kit

# Check ownership
sudo chown $USER:$USER /usr/local/bin/starknode-kit
```

#### Go Version Issues

**Problem**: Go version too old

**Solutions**:
```bash
# Check current version
go version

# Update Go (Ubuntu/Debian)
sudo apt update
sudo apt install golang-go

# Update Go (macOS)
brew install go

# Update Go (manual)
# Visit https://go.dev/dl/ for latest version
```

---

### Configuration Issues

#### Invalid Configuration

**Problem**: Configuration validation fails

**Solutions**:
```bash
# Validate configuration
starknode-kit config --validate

# Reset to defaults
rm ~/.starknode-kit/config.yaml
starknode-kit init

# Check YAML syntax
yamllint ~/.starknode-kit/config.yaml
```

#### Missing Configuration File

**Problem**: Configuration file not found

**Solutions**:
```bash
# Create default configuration
starknode-kit init

# Check file location
ls -la ~/.starknode-kit/

# Set custom path
export STARKNODE_CONFIG_PATH=/custom/path/config.yaml
```

---

### Client Installation Issues

#### Download Failures

**Problem**: Client binaries fail to download

**Solutions**:
```bash
# Check internet connection
ping google.com

# Check available disk space
df -h

# Check GitHub access
curl -I https://github.com

# Manual download
wget https://github.com/ethereum/go-ethereum/releases/download/v1.13.5/geth-linux-amd64-1.13.5-916d9836.tar.gz
```

#### Installation Permissions

**Problem**: Cannot install client binaries

**Solutions**:
```bash
# Check directory permissions
ls -la ~/.starknode-kit/

# Fix permissions
chmod 755 ~/.starknode-kit/
chmod 755 ~/.starknode-kit/bin/

# Create directories if missing
mkdir -p ~/.starknode-kit/bin/
```

---

### Client Runtime Issues

#### Port Conflicts

**Problem**: Clients cannot bind to required ports

**Solutions**:
```bash
# Check port usage
netstat -tulpn | grep :8545  # Geth
netstat -tulpn | grep :9000  # Lighthouse
netstat -tulpn | grep :6060  # Juno

# Kill conflicting processes
sudo kill -9 $(lsof -t -i:8545)

# Use different ports
starknode-kit config el client=geth port=8547,8548
```

#### Sync Stalls

**Problem**: Blockchain sync stops or becomes very slow

**Solutions**:
```bash
# Check client status
starknode-kit monitor

# Restart clients
starknode-kit stop
starknode-kit start

# Check disk space
df -h

# Check network connectivity
ping 8.8.8.8

# Use checkpoint sync
starknode-kit config cl client=lighthouse consensus_checkpoint="https://mainnet.checkpoint.sigp.io"
```

#### High Resource Usage

**Problem**: Clients consume too much CPU/memory

**Solutions**:
```bash
# Check resource usage
htop
iotop

# Reduce client resources
starknode-kit config el client=geth additional_args="--cache,2048,--maxpeers,25"

# Use lighter execution type
starknode-kit config el client=geth execution_type=light
```

---

### Network Issues

#### Connection Problems

**Problem**: Cannot connect to network

**Solutions**:
```bash
# Check network connectivity
ping google.com

# Check DNS resolution
nslookup google.com

# Check firewall
sudo ufw status

# Test specific ports
telnet 8.8.8.8 53
```

#### Low Peer Count

**Problem**: Few or no peer connections

**Solutions**:
```bash
# Check peer count
starknode-kit monitor

# Increase max peers
starknode-kit config el client=geth additional_args="--maxpeers,100"

# Check network configuration
starknode-kit config -n mainnet

# Restart clients
starknode-kit stop && starknode-kit start
```

---

### Monitoring Issues

#### Dashboard Not Loading

**Problem**: Monitor dashboard fails to start

**Solutions**:
```bash
# Check if clients are running
starknode-kit config

# Restart monitoring
starknode-kit monitor --restart

# Check for port conflicts
netstat -tulpn | grep :9090

# Use basic monitoring
starknode-kit monitor --basic
```

#### Missing Metrics

**Problem**: Some metrics not showing in dashboard

**Solutions**:
```bash
# Verify client status
starknode-kit config

# Check client logs
tail -f ~/.starknode-kit/logs/*.log

# Restart clients
starknode-kit stop && starknode-kit start

# Check monitoring configuration
starknode-kit monitor --interval 10s
```

---

## Diagnostic Commands

### System Information

```bash
# Check system resources
free -h
df -h
lscpu

# Check network interfaces
ip addr show
ip route show

# Check running processes
ps aux | grep -E "(geth|lighthouse|juno)"
```

### Client Status

```bash
# Check client processes
pgrep -f geth
pgrep -f lighthouse
pgrep -f juno

# Check client logs
tail -f ~/.starknode-kit/logs/geth.log
tail -f ~/.starknode-kit/logs/lighthouse.log
tail -f ~/.starknode-kit/logs/juno.log

# Check client configuration
starknode-kit config
```

### Network Diagnostics

```bash
# Test connectivity
ping -c 4 8.8.8.8
ping -c 4 google.com

# Check DNS
nslookup google.com
dig google.com

# Test specific ports
nc -zv 8.8.8.8 53
nc -zv localhost 8545
```

---

## Log Analysis

### Common Log Patterns

#### Geth Logs

**Sync Issues**:
```bash
grep -i "sync" ~/.starknode-kit/logs/geth.log
grep -i "peer" ~/.starknode-kit/logs/geth.log
```

**Memory Issues**:
```bash
grep -i "memory" ~/.starknode-kit/logs/geth.log
grep -i "cache" ~/.starknode-kit/logs/geth.log
```

#### Lighthouse Logs

**Consensus Issues**:
```bash
grep -i "consensus" ~/.starknode-kit/logs/lighthouse.log
grep -i "beacon" ~/.starknode-kit/logs/lighthouse.log
```

**Network Issues**:
```bash
grep -i "network" ~/.starknode-kit/logs/lighthouse.log
grep -i "connection" ~/.starknode-kit/logs/lighthouse.log
```

#### Juno Logs

**L1 Verification Issues**:
```bash
grep -i "l1" ~/.starknode-kit/logs/juno.log
grep -i "ethereum" ~/.starknode-kit/logs/juno.log
```

**Sync Issues**:
```bash
grep -i "sync" ~/.starknode-kit/logs/juno.log
grep -i "block" ~/.starknode-kit/logs/juno.log
```

---

## Performance Issues

### Slow Sync

**Causes**:
- Insufficient disk I/O
- Low memory
- Poor network connection
- High CPU usage

**Solutions**:
```bash
# Check disk I/O
iotop
iostat -x 1

# Check memory usage
free -h
htop

# Optimize client settings
starknode-kit config el client=geth additional_args="--cache,4096,--maxpeers,50"
```

### High Memory Usage

**Causes**:
- Large cache settings
- Memory leaks
- Insufficient swap

**Solutions**:
```bash
# Check memory usage
free -h
htop

# Reduce cache size
starknode-kit config el client=geth additional_args="--cache,2048"

# Add swap space
sudo fallocate -l 8G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
```

### High CPU Usage

**Causes**:
- Intensive sync process
- High peer count
- Resource contention

**Solutions**:
```bash
# Check CPU usage
htop
top

# Reduce peer count
starknode-kit config el client=geth additional_args="--maxpeers,25"

# Use lighter execution type
starknode-kit config el client=geth execution_type=light
```

---

## Recovery Procedures

### Complete Reset

```bash
# Stop all clients
starknode-kit stop

# Remove configuration
rm -rf ~/.starknode-kit/

# Reinstall
curl -sSL https://raw.githubusercontent.com/thebuidl-grid/starknode-kit/main/install.sh | bash

# Reinitialize
starknode-kit init
```

### Client Reset

```bash
# Stop specific client
starknode-kit stop

# Remove client
starknode-kit remove --execution_client geth

# Reinstall client
starknode-kit add --execution_client geth

# Restart
starknode-kit start
```

### Database Reset

```bash
# Stop clients
starknode-kit stop

# Remove database (WARNING: This will require full resync)
rm -rf ~/.starknode-kit/data/geth/
rm -rf ~/.starknode-kit/data/lighthouse/
rm -rf ~/.starknode-kit/data/juno/

# Restart clients
starknode-kit start
```

---

## Getting Help

### Before Asking for Help

1. **Check logs**: Review client logs for error messages
2. **Verify configuration**: Ensure configuration is valid
3. **Check system resources**: Verify adequate disk space and memory
4. **Test network**: Ensure stable internet connection
5. **Try restart**: Restart clients and check if issue persists

### Information to Provide

When seeking help, include:

- **Operating system**: Linux distribution, macOS version, etc.
- **Hardware specs**: CPU, RAM, storage
- **Starknode Kit version**: `starknode-kit --version`
- **Configuration**: `starknode-kit config`
- **Error messages**: Full error output
- **Logs**: Relevant log excerpts
- **Steps to reproduce**: What you did before the issue occurred

### Support Channels

- **Telegram**: [Join our community](https://t.me/+SCPbza9fk8dkYWI0)
- **GitHub Issues**: [Report bugs and request features](https://github.com/thebuidl-grid/starknode-kit/issues)
- **Documentation**: This troubleshooting guide

### Community Resources

- **Discord**: Real-time community support
- **Reddit**: r/ethereum, r/starknet
- **Stack Overflow**: Tag questions with `starknode-kit`
- **GitHub Discussions**: Community discussions and Q&A
