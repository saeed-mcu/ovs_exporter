# Migration Guide

## Migrating from greenpau/ovs_exporter to Liquescent-Development/ovs_exporter

This guide helps you migrate from the original ovs_exporter to this enhanced fork.

## Key Improvements in v2.0.0

- **Go 1.24** - Modernized to latest Go version
- **Updated Dependencies** - All dependencies updated to latest versions with security patches
- **PMD/DPDK Metrics** - Comprehensive performance metrics for DPDK deployments
- **Security Updates** - Fixed vulnerabilities in protobuf and other dependencies
- **Enhanced Coverage** - All metrics from Red Hat's OVS observability blog post

## Migration Steps

### 1. Update Go Module References

If you're importing this module in your Go code:

```go
// Old
import "github.com/greenpau/ovs_exporter/pkg/ovs_exporter"

// New
import "github.com/Liquescent-Development/ovs_exporter/pkg/ovs_exporter"
```

### 2. Update go.mod

```bash
go get github.com/Liquescent-Development/ovs_exporter@v2.0.0
```

Or in your go.mod:

```go
require (
    github.com/Liquescent-Development/ovs_exporter v2.0.0
)
```

### 3. Update Docker/Kubernetes Deployments

If using Docker:

```dockerfile
# Old
FROM greenpau/ovs-exporter:latest

# New - Build from source or wait for official images
FROM golang:1.24 as builder
RUN git clone https://github.com/Liquescent-Development/ovs_exporter.git
WORKDIR /ovs_exporter
RUN make
```

### 4. Update Systemd Service

The binary name remains `ovs-exporter`, but if you have hardcoded paths:

```bash
# Download new version
wget https://github.com/Liquescent-Development/ovs_exporter/releases/download/v2.0.0/ovs-exporter-2.0.0.linux-amd64.tar.gz

# Install
tar xvzf ovs-exporter-2.0.0.linux-amd64.tar.gz
cd ovs-exporter-2.0.0.linux-amd64
sudo ./install.sh
```

## New Metrics Available

After migration, you'll have access to new PMD/DPDK metrics:

- `ovs_pmd_cycles_per_iteration`
- `ovs_pmd_packets_per_iteration`
- `ovs_pmd_cycles_per_packet`
- `ovs_pmd_packets_per_batch`
- `ovs_vhost_tx_retries_total`
- `ovs_vhost_tx_contention_total`
- And many more (see [PMD_METRICS.md](PMD_METRICS.md))

## Compatibility

- **Backward Compatible**: All existing metrics remain unchanged
- **Drop-in Replacement**: No configuration changes required
- **Graceful PMD Handling**: PMD metrics collection automatically skips on non-DPDK systems

## Prometheus Configuration

No changes needed to your Prometheus configuration. The exporter uses the same:
- Port: 9475
- Path: /metrics
- Metric namespace: `ovs_`

## Grafana Dashboards

Existing dashboards will continue to work. Consider adding new panels for PMD metrics:

```promql
# Example: PMD efficiency
rate(ovs_pmd_iterations_total[5m]) * ovs_pmd_packets_per_iteration

# Example: vHost contention
rate(ovs_vhost_tx_contention_total[5m])
```

## Rollback Plan

If you need to rollback:

1. Stop the new exporter
2. Install the old version
3. Update go.mod to use `github.com/greenpau/ovs_exporter`
4. Note: You'll lose access to PMD metrics

## Support

- Issues: https://github.com/Liquescent-Development/ovs_exporter/issues
- Original project: https://github.com/greenpau/ovs_exporter

## Version Mapping

| Original Version | Fork Version | Notes |
|-----------------|--------------|-------|
| v1.0.7 | v2.0.0 | Major upgrade with PMD metrics and Go 1.24 |