# Monitoring

Starknode Kit includes a built-in monitoring dashboard that provides real-time insights into your node's health, performance, and synchronization status.

## Built-in Monitoring Dashboard

### Launch the Dashboard

```bash
starknode-kit monitor
```

This opens an interactive terminal-based dashboard showing:

- **Client Status**: Running/stopped status of all clients
- **Sync Progress**: Current synchronization status
- **Resource Usage**: CPU, memory, and disk usage
- **Network Stats**: Peer connections and bandwidth
- **Logs**: Real-time log output from clients

### Dashboard Features

#### Real-time Updates
- Automatic refresh every few seconds
- Live log streaming
- Status change notifications

#### Multi-client View
- Execution client status
- Consensus client status
- Starknet client status
- Combined system overview

#### Resource Monitoring
- CPU usage per client
- Memory consumption
- Disk I/O statistics
- Network bandwidth

## Monitoring Components

### Client Status Monitoring

**Execution Client (Geth/Reth)**:
- Sync status and progress
- Peer connections
- Block height
- Transaction pool status

**Consensus Client (Lighthouse/Prysm)**:
- Beacon chain sync
- Validator status (if applicable)
- Attestation performance
- Finality status

**Starknet Client (Juno)**:
- L2 sync progress
- L1 verification status
- Transaction processing
- State root updates

### System Resource Monitoring

**CPU Usage**:
- Per-client CPU consumption
- System-wide CPU usage
- Load average

**Memory Usage**:
- RAM consumption per client
- Available memory
- Swap usage

**Disk Usage**:
- Database size
- Available disk space
- I/O operations per second

**Network Monitoring**:
- Inbound/outbound bandwidth
- Peer connections
- Latency measurements

## Advanced Monitoring

### Custom Monitoring Scripts

Create custom monitoring scripts for specific needs:

```bash
#!/bin/bash
# monitor_health.sh

# Check if clients are running
if ! pgrep -f "geth" > /dev/null; then
    echo "Geth is not running!"
    exit 1
fi

if ! pgrep -f "lighthouse" > /dev/null; then
    echo "Lighthouse is not running!"
    exit 1
fi

# Check sync status
GETH_SYNC=$(curl -s -X POST -H "Content-Type: application/json" \
    --data '{"jsonrpc":"2.0","method":"eth_syncing","params":[],"id":1}' \
    http://localhost:8545 | jq -r '.result')

if [ "$GETH_SYNC" != "false" ]; then
    echo "Geth is still syncing"
fi

echo "All systems healthy"
```

### Log Monitoring

Monitor specific log patterns:

```bash
# Monitor Geth logs
tail -f ~/.starknode-kit/logs/geth.log | grep -E "(ERROR|WARN|FATAL)"

# Monitor Lighthouse logs
tail -f ~/.starknode-kit/logs/lighthouse.log | grep -E "(ERROR|WARN|FATAL)"

# Monitor Juno logs
tail -f ~/.starknode-kit/logs/juno.log | grep -E "(ERROR|WARN|FATAL)"
```

### Performance Metrics

Track key performance indicators:

```bash
# Check sync speed
starknode-kit monitor --metrics sync-speed

# Check peer count
starknode-kit monitor --metrics peers

# Check resource usage
starknode-kit monitor --metrics resources
```

## External Monitoring Tools

### Prometheus Integration

Set up Prometheus metrics collection:

```yaml
# prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'starknode-kit'
    static_configs:
      - targets: ['localhost:9090']
    metrics_path: /metrics
    scrape_interval: 5s
```

### Grafana Dashboards

Create Grafana dashboards for visualization:

**Key Metrics to Track**:
- Sync progress percentage
- Block height over time
- Peer connection count
- Resource usage trends
- Error rates

### Alerting

Set up alerts for critical issues:

**Critical Alerts**:
- Client process down
- Sync stalled for > 1 hour
- Disk space < 10%
- Memory usage > 90%
- High error rate

**Warning Alerts**:
- Sync progress slow
- Low peer count
- High CPU usage
- Network connectivity issues

## Monitoring Best Practices

### Regular Health Checks

**Daily Checks**:
- Verify all clients are running
- Check sync progress
- Review error logs
- Monitor resource usage

**Weekly Checks**:
- Review performance trends
- Check disk space usage
- Update monitoring configurations
- Test alerting systems

### Log Management

**Log Rotation**:
```bash
# Set up log rotation
sudo logrotate -f /etc/logrotate.d/starknode-kit
```

**Log Analysis**:
```bash
# Analyze error patterns
grep -E "(ERROR|WARN)" ~/.starknode-kit/logs/*.log | \
    awk '{print $1}' | sort | uniq -c | sort -nr
```

### Performance Optimization

**Based on Monitoring Data**:
- Adjust client configurations
- Optimize resource allocation
- Tune network settings
- Update hardware if needed

## Troubleshooting Monitoring Issues

### Dashboard Not Loading

**Problem**: Monitor dashboard fails to start
**Solutions**:
```bash
# Check if clients are configured
starknode-kit config

# Restart monitoring
starknode-kit monitor --restart

# Check for port conflicts
netstat -tulpn | grep :9090
```

### Missing Metrics

**Problem**: Some metrics not showing
**Solutions**:
```bash
# Verify client status
starknode-kit config

# Check client logs
tail -f ~/.starknode-kit/logs/*.log

# Restart clients
starknode-kit stop && starknode-kit start
```

### High Resource Usage

**Problem**: Monitoring consumes too many resources
**Solutions**:
```bash
# Reduce monitoring frequency
starknode-kit monitor --interval 30s

# Disable detailed metrics
starknode-kit monitor --basic

# Use external monitoring instead
```

## Monitoring Configuration

### Customize Monitoring Settings

```yaml
# ~/.starknode-kit/monitoring.yaml
monitoring:
  refresh_interval: 5s
  log_level: info
  metrics_enabled: true
  alerts_enabled: true
  
  clients:
    geth:
      metrics_port: 6060
      log_level: info
    lighthouse:
      metrics_port: 5054
      log_level: info
    juno:
      metrics_port: 9090
      log_level: info
```

### Environment Variables

```bash
# Set monitoring environment variables
export STARKNODE_MONITOR_INTERVAL=10s
export STARKNODE_MONITOR_LOG_LEVEL=debug
export STARKNODE_MONITOR_METRICS=true
```

## Integration with External Tools

### System Monitoring

**htop/btop**:
```bash
# Monitor system resources
htop
btop
```

**iotop**:
```bash
# Monitor disk I/O
sudo iotop
```

**nethogs**:
```bash
# Monitor network usage
sudo nethogs
```

### Log Aggregation

**ELK Stack**:
- Elasticsearch for log storage
- Logstash for log processing
- Kibana for visualization

**Fluentd**:
- Lightweight log collector
- Easy integration with Starknode Kit

## Getting Help

For monitoring issues:

- Check the [Troubleshooting Guide](troubleshooting.md)
- Review client-specific monitoring docs
- Join our [Telegram community](https://t.me/+SCPbza9fk8dkYWI0)
- Open an issue on [GitHub](https://github.com/thebuidl-grid/starknode-kit/issues)
