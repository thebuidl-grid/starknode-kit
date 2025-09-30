# Updates and Maintenance

This guide covers how to keep your Starknode Kit installation and clients up to date, along with best practices for maintenance.

## Client Updates

### Check for Updates

```bash
# Check for updates to all clients
starknode-kit update

# Check for updates to specific client
starknode-kit update --execution_client geth
starknode-kit update --consensus_client lighthouse
starknode-kit update --starknet_client juno
```

### Update Process

When you run `starknode-kit update`, the following happens:

1. **Check versions**: Compare current vs. latest available versions
2. **Download updates**: Download new client binaries if available
3. **Backup current**: Backup existing binaries
4. **Install updates**: Replace old binaries with new ones
5. **Validate installation**: Verify new binaries work correctly

### Automatic Updates

Enable automatic update checking:

```bash
# Set up cron job for daily updates
echo "0 2 * * * /usr/local/bin/starknode-kit update" | crontab -

# Or weekly updates
echo "0 2 * * 0 /usr/local/bin/starknode-kit update" | crontab -
```

---

## Update Strategies

### Rolling Updates

Update clients one at a time to minimize downtime:

```bash
# 1. Update execution client
starknode-kit stop
starknode-kit update --execution_client geth
starknode-kit start

# 2. Wait for sync, then update consensus client
starknode-kit stop
starknode-kit update --consensus_client lighthouse
starknode-kit start

# 3. Update Starknet client
starknode-kit stop
starknode-kit update --starknet_client juno
starknode-kit run juno
```

### Maintenance Windows

Schedule updates during low-activity periods:

```bash
# Create maintenance script
cat > ~/maintenance.sh << 'EOF'
#!/bin/bash
echo "Starting maintenance window..."

# Stop clients
starknode-kit stop

# Update all clients
starknode-kit update

# Clean up old logs
find ~/.starknode-kit/logs -name "*.log" -mtime +30 -delete

# Restart clients
starknode-kit start

echo "Maintenance complete"
EOF

chmod +x ~/maintenance.sh
```

---

## Version Management

### Check Current Versions

```bash
# Check Starknode Kit version
starknode-kit --version

# Check client versions
~/.starknode-kit/bin/geth version
~/.starknode-kit/bin/lighthouse --version
~/.starknode-kit/bin/juno --version
```

### Version Compatibility

| Starknode Kit | Geth | Lighthouse | Prysm | Juno |
|---------------|------|------------|-------|------|
| 1.0.0 | 1.13.5+ | 4.5.0+ | 4.0.0+ | 0.14.0+ |
| 1.1.0 | 1.13.8+ | 4.6.0+ | 4.1.0+ | 0.14.5+ |
| 1.2.0 | 1.14.0+ | 4.7.0+ | 4.2.0+ | 0.15.0+ |

### Downgrade if Needed

```bash
# Stop clients
starknode-kit stop

# Remove current version
rm ~/.starknode-kit/bin/geth

# Download specific version
wget https://github.com/ethereum/go-ethereum/releases/download/v1.13.5/geth-linux-amd64-1.13.5-916d9836.tar.gz
tar -xzf geth-linux-amd64-1.13.5-916d9836.tar.gz
mv geth ~/.starknode-kit/bin/

# Restart clients
starknode-kit start
```

---

## Maintenance Tasks

### Regular Maintenance

#### Daily Tasks

```bash
# Check client status
starknode-kit monitor

# Check disk space
df -h

# Check logs for errors
grep -i error ~/.starknode-kit/logs/*.log
```

#### Weekly Tasks

```bash
# Update clients
starknode-kit update

# Clean old logs
find ~/.starknode-kit/logs -name "*.log" -mtime +7 -delete

# Check system resources
htop
iotop
```

#### Monthly Tasks

```bash
# Full system update
sudo apt update && sudo apt upgrade

# Clean package cache
sudo apt autoremove
sudo apt autoclean

# Check disk health
sudo smartctl -a /dev/sda
```

### Log Management

#### Log Rotation

```bash
# Set up logrotate
sudo tee /etc/logrotate.d/starknode-kit << 'EOF'
/home/*/.starknode-kit/logs/*.log {
    daily
    missingok
    rotate 7
    compress
    delaycompress
    notifempty
    create 644 $USER $USER
}
EOF
```

#### Log Analysis

```bash
# Analyze error patterns
grep -E "(ERROR|WARN|FATAL)" ~/.starknode-kit/logs/*.log | \
    awk '{print $1}' | sort | uniq -c | sort -nr

# Check sync performance
grep -i "sync" ~/.starknode-kit/logs/geth.log | tail -20

# Monitor peer connections
grep -i "peer" ~/.starknode-kit/logs/*.log | tail -20
```

---

## Backup and Recovery

### Configuration Backup

```bash
# Backup configuration
cp ~/.starknode-kit/config.yaml ~/.starknode-kit/config.yaml.backup

# Backup to cloud storage
rclone copy ~/.starknode-kit/config.yaml remote:backups/

# Version control
cd ~/.starknode-kit/
git init
git add config.yaml
git commit -m "Initial configuration"
```

### Database Backup

```bash
# Create backup script
cat > ~/backup_databases.sh << 'EOF'
#!/bin/bash
BACKUP_DIR="/backup/starknode-kit/$(date +%Y%m%d)"
mkdir -p "$BACKUP_DIR"

# Stop clients
starknode-kit stop

# Backup databases
tar -czf "$BACKUP_DIR/geth_data.tar.gz" ~/.starknode-kit/data/geth/
tar -czf "$BACKUP_DIR/lighthouse_data.tar.gz" ~/.starknode-kit/data/lighthouse/
tar -czf "$BACKUP_DIR/juno_data.tar.gz" ~/.starknode-kit/data/juno/

# Restart clients
starknode-kit start

echo "Backup completed: $BACKUP_DIR"
EOF

chmod +x ~/backup_databases.sh
```

### Recovery Procedures

```bash
# Restore configuration
cp ~/.starknode-kit/config.yaml.backup ~/.starknode-kit/config.yaml

# Restore database
starknode-kit stop
tar -xzf /backup/starknode-kit/20240101/geth_data.tar.gz -C ~/.starknode-kit/data/
starknode-kit start
```

---

## Performance Monitoring

### System Metrics

```bash
# Monitor system resources
htop
iotop
nethogs

# Check disk usage
du -sh ~/.starknode-kit/data/*

# Monitor network
iftop
```

### Client Metrics

```bash
# Monitor sync progress
starknode-kit monitor

# Check peer count
curl -X POST -H "Content-Type: application/json" \
    --data '{"jsonrpc":"2.0","method":"net_peerCount","params":[],"id":1}' \
    http://localhost:8545

# Check block height
curl -X POST -H "Content-Type: application/json" \
    --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
    http://localhost:8545
```

---

## Security Updates

### System Security

```bash
# Update system packages
sudo apt update && sudo apt upgrade

# Check for security updates
sudo apt list --upgradable

# Update kernel if needed
sudo apt install linux-image-generic
```

### Client Security

```bash
# Update clients regularly
starknode-kit update

# Check for security advisories
# Monitor client GitHub repositories for security releases

# Use latest stable versions
starknode-kit config --validate
```

---

## Troubleshooting Updates

### Update Failures

**Problem**: Update process fails

**Solutions**:
```bash
# Check internet connection
ping google.com

# Check disk space
df -h

# Check permissions
ls -la ~/.starknode-kit/bin/

# Manual update
wget https://github.com/ethereum/go-ethereum/releases/download/v1.14.0/geth-linux-amd64-1.14.0.tar.gz
tar -xzf geth-linux-amd64-1.14.0.tar.gz
mv geth ~/.starknode-kit/bin/
```

### Version Conflicts

**Problem**: New version causes issues

**Solutions**:
```bash
# Rollback to previous version
starknode-kit stop
rm ~/.starknode-kit/bin/geth
# Restore previous version from backup

# Check compatibility
starknode-kit config --validate

# Report issue
# Open GitHub issue with version information
```

### Sync Issues After Update

**Problem**: Sync problems after updating

**Solutions**:
```bash
# Check client logs
tail -f ~/.starknode-kit/logs/geth.log

# Restart clients
starknode-kit stop && starknode-kit start

# Use checkpoint sync
starknode-kit config cl client=lighthouse consensus_checkpoint="https://mainnet.checkpoint.sigp.io"
```

---

## Automation

### Update Automation

```bash
# Create update script
cat > ~/auto_update.sh << 'EOF'
#!/bin/bash
LOG_FILE="/var/log/starknode-kit-updates.log"

echo "$(date): Starting automatic update" >> "$LOG_FILE"

# Check for updates
if starknode-kit update --check-only; then
    echo "$(date): Updates available, proceeding" >> "$LOG_FILE"
    
    # Stop clients
    starknode-kit stop
    
    # Update clients
    starknode-kit update
    
    # Restart clients
    starknode-kit start
    
    echo "$(date): Update completed successfully" >> "$LOG_FILE"
else
    echo "$(date): No updates available" >> "$LOG_FILE"
fi
EOF

chmod +x ~/auto_update.sh

# Schedule weekly updates
echo "0 3 * * 0 /home/$USER/auto_update.sh" | crontab -
```

### Monitoring Automation

```bash
# Create health check script
cat > ~/health_check.sh << 'EOF'
#!/bin/bash
ALERT_EMAIL="admin@example.com"

# Check if clients are running
if ! pgrep -f "geth" > /dev/null; then
    echo "Geth is not running!" | mail -s "Starknode Kit Alert" "$ALERT_EMAIL"
fi

if ! pgrep -f "lighthouse" > /dev/null; then
    echo "Lighthouse is not running!" | mail -s "Starknode Kit Alert" "$ALERT_EMAIL"
fi

# Check disk space
DISK_USAGE=$(df / | awk 'NR==2 {print $5}' | sed 's/%//')
if [ "$DISK_USAGE" -gt 90 ]; then
    echo "Disk usage is ${DISK_USAGE}%" | mail -s "Starknode Kit Alert" "$ALERT_EMAIL"
fi
EOF

chmod +x ~/health_check.sh

# Schedule hourly health checks
echo "0 * * * * /home/$USER/health_check.sh" | crontab -
```

---

## Best Practices

### Update Strategy

1. **Test on testnets first**: Always test updates on testnets before mainnet
2. **Schedule maintenance windows**: Update during low-activity periods
3. **Keep backups**: Always backup before major updates
4. **Monitor after updates**: Watch for issues after updating
5. **Document changes**: Keep track of what was updated and when

### Maintenance Schedule

- **Daily**: Check client status and logs
- **Weekly**: Update clients and clean logs
- **Monthly**: Full system maintenance and security updates
- **Quarterly**: Review and optimize configurations

### Monitoring

- Set up alerts for critical issues
- Monitor resource usage trends
- Track sync performance
- Watch for security updates

---

## Getting Help

For update and maintenance issues:

- Check the [Troubleshooting Guide](troubleshooting.md)
- Review client-specific update documentation
- Join our [Telegram community](https://t.me/+SCPbza9fk8dkYWI0)
- Open an issue on [GitHub](https://github.com/thebuidl-grid/starknode-kit/issues)
