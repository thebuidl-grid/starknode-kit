# Hardware Requirements

This guide outlines the hardware requirements for running Ethereum and Starknet nodes with Starknode Kit. Requirements vary depending on your use case, from casual development to production-grade operations.

## Minimum Requirements

### Development/Testing Environment

- **CPU**: 4+ cores (Intel i5/AMD Ryzen 5 or better)
- **RAM**: 16 GB
- **Storage**: 1 TB NVMe SSD
- **Network**: Stable broadband connection (10+ Mbps)

### Production Environment

- **CPU**: 8+ cores (Intel i7/AMD Ryzen 7 or better)
- **RAM**: 32+ GB
- **Storage**: 2+ TB NVMe SSD with DRAM cache
- **Network**: Stable broadband connection (50+ Mbps)

## Detailed Hardware Specifications

### CPU Requirements

**Recommended**: Intel i7/AMD Ryzen 7 or better
- Node operation doesn't require heavy CPU power
- The BG Client has run well on both i3 and i5 models
- **Avoid**: Celeron processors due to limitations

**Why CPU matters**:
- Consensus client operations
- Block processing and validation
- Network protocol handling

### Memory (RAM) Requirements

**Minimum**: 16 GB
**Recommended**: 32+ GB
**Production**: 64+ GB

**Memory usage breakdown**:
- Base system: ~4 GB
- Ethereum execution client: ~8-16 GB
- Consensus client: ~4-8 GB
- Starknet client (Juno): ~4-8 GB
- Operating system overhead: ~4 GB

### Storage Requirements

**Critical**: Use NVMe SSD with DRAM cache
**Avoid**: QLC (Quad-Level Cell) NAND architecture

#### Storage Size Requirements

| Network | Execution Client | Consensus Client | Total (Approx.) |
|---------|------------------|------------------|-----------------|
| Mainnet | 1.5 TB | 500 GB | 2+ TB |
| Sepolia | 200 GB | 50 GB | 300+ GB |
| Goerli | 300 GB | 100 GB | 500+ GB |

#### Recommended SSDs

See this [SSD List Gist](https://gist.github.com/bkase/fab02c5b3c404e9ef8e5c2071ac1558c) for tested options.

**Key specifications to look for**:
- NVMe interface
- DRAM cache
- High endurance rating
- Good random read/write performance

### Network Requirements

**Bandwidth**:
- Development: 10+ Mbps
- Production: 50+ Mbps
- High-performance: 100+ Mbps

**Latency**:
- Low latency connection preferred
- Stable connection more important than speed

**Data usage**:
- Initial sync: 1-2 TB
- Daily operation: 10-50 GB

## Hardware Recommendations by Use Case

### Home Development Setup

```
CPU: Intel i5-12400 / AMD Ryzen 5 5600X
RAM: 32 GB DDR4
Storage: 2 TB NVMe SSD (Samsung 980 Pro, WD Black SN850)
Network: 50+ Mbps broadband
```

### Production Validator

```
CPU: Intel i7-12700K / AMD Ryzen 7 5800X
RAM: 64 GB DDR4
Storage: 4 TB NVMe SSD (Samsung 980 Pro, WD Black SN850)
Network: 100+ Mbps dedicated connection
```

### Enterprise/High-Performance

```
CPU: Intel i9-12900K / AMD Ryzen 9 5900X
RAM: 128 GB DDR4
Storage: 8 TB NVMe SSD (Enterprise grade)
Network: 1 Gbps dedicated connection
```

## Cloud Provider Recommendations

### AWS
- **Instance**: c6i.2xlarge or larger
- **Storage**: gp3 or io2 volumes
- **Network**: Enhanced networking enabled

### Google Cloud
- **Instance**: c2-standard-8 or larger
- **Storage**: SSD persistent disks
- **Network**: Premium tier

### Azure
- **Instance**: Standard_D8s_v3 or larger
- **Storage**: Premium SSD
- **Network**: Accelerated networking

## Performance Optimization Tips

### Storage Optimization
- Use separate drives for different clients
- Enable TRIM for SSD maintenance
- Monitor disk health regularly

### Memory Optimization
- Enable swap space (8-16 GB)
- Monitor memory usage patterns
- Consider memory overcommit settings

### Network Optimization
- Use wired connections when possible
- Configure QoS for node traffic
- Monitor bandwidth usage

## Monitoring Hardware Health

### Key Metrics to Monitor
- CPU usage and temperature
- Memory usage and swap usage
- Disk I/O and available space
- Network bandwidth and latency

### Recommended Tools
- `htop` for system monitoring
- `iotop` for disk I/O monitoring
- `nethogs` for network monitoring
- `smartctl` for SSD health

## Troubleshooting Hardware Issues

### Common Problems

**High CPU Usage**
- Check for stuck processes
- Verify client configurations
- Monitor system logs

**Memory Issues**
- Check for memory leaks
- Adjust client memory limits
- Consider adding swap space

**Storage Problems**
- Monitor disk space
- Check for disk errors
- Verify SSD health

**Network Issues**
- Test connection stability
- Check firewall settings
- Verify port configurations

## Getting Help

For hardware-specific questions:

- Check our [Troubleshooting Guide](../operations/troubleshooting.md)
- Join our [Telegram community](https://t.me/+SCPbza9fk8dkYWI0)
- Review the [Rocket Pool Hardware Guide](https://docs.rocketpool.net/guides/node/hardware.html) for additional insights
