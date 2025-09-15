# OVS Exporter Grafana Dashboards

This directory contains Grafana dashboards designed to visualize metrics from the OVS Exporter. These dashboards showcase the ability to monitor OVS performance at both executive and operational levels, with special support for multi-tenant environments.

## Dashboard Overview

### 1. Executive SLA/SLO Dashboard (`ovs-executive-sla-slo.json`)

**Purpose**: High-level view for executives and SLA/SLO monitoring with per-tenant performance tracking.

**Key Features**:
- **Service Level Overview**: Real-time system status, success rates, and CPU utilization trends
- **Per-Tenant Metrics**: Packet rates and bandwidth consumption grouped by tenant ID (via external_ids)
- **Tenant Performance Table**: Comprehensive summary showing packet totals, drop rates, and error rates per tenant
- **SLA/SLO Indicators**:
  - Cache performance metrics (EMC, SMC, Megaflow hit rates)
  - Drop reason distribution
  - Packet processing latency percentiles (P95, P99)

**Use Cases**:
- Executive reporting on customer experience
- SLA compliance monitoring
- Multi-tenant performance comparison
- Quick identification of tenant-specific issues

### 2. Operational Dashboard (`ovs-operational.json`)

**Purpose**: Detailed technical view for operators and engineers to monitor and troubleshoot OVS performance.

**Key Features**:
- **System Health**: OVS status, interface counts, flow counts, failed request rates
- **Memory Monitoring**: Real-time memory usage and log file sizes
- **PMD Performance**:
  - CPU utilization per PMD thread
  - Cycles per packet metrics
  - Packet processing rates
  - Flow lookup statistics
- **Interface Metrics**:
  - Bandwidth and packet rates per interface
  - Error and drop statistics
  - Interface utilization percentages
- **Flow Cache Performance**: Detailed hit ratios for all cache levels
- **Drop Analysis**: Time-series view of drops by reason

**Use Cases**:
- Performance troubleshooting
- Capacity planning
- Bottleneck identification
- Detailed system monitoring

## Installation

### Prerequisites

1. **Grafana**: Version 9.0 or later recommended
2. **Prometheus**: Configured to scrape OVS exporter metrics
3. **OVS Exporter**: Running and exposing metrics (default port 9475)

### Import Process

1. **Via Grafana UI**:
   ```
   1. Navigate to Dashboards → Import
   2. Upload JSON file or paste JSON content
   3. Select your Prometheus datasource
   4. Click "Import"
   ```

2. **Via Grafana API**:
   ```bash
   curl -X POST http://admin:admin@localhost:3000/api/dashboards/db \
     -H "Content-Type: application/json" \
     -d @ovs-executive-sla-slo.json
   ```

3. **Via Provisioning** (recommended for production):
   ```yaml
   # /etc/grafana/provisioning/dashboards/ovs.yaml
   apiVersion: 1
   providers:
     - name: 'OVS Dashboards'
       orgId: 1
       folder: 'OVS'
       type: file
       disableDeletion: false
       updateIntervalSeconds: 10
       allowUiUpdates: true
       options:
         path: /var/lib/grafana/dashboards/ovs
   ```

## Configuration

### Tenant ID Mapping

For proper tenant grouping, ensure your OVS interfaces have external_ids configured:

```bash
# Set tenant ID for an interface
ovs-vsctl set Interface <interface-name> external_ids:tenant_id="customer-123"

# Verify configuration
ovs-vsctl get Interface <interface-name> external_ids
```

The dashboards will automatically group metrics by these tenant IDs, enabling per-customer performance tracking.

### Dashboard Variables

Both dashboards include configurable variables:

- **datasource**: Select your Prometheus instance
- **system_id**: Filter by OVS system ID (auto-populated)
- **interface**: (Operational dashboard only) Filter by specific interface

### Refresh Intervals

- Executive Dashboard: 30 seconds (configurable)
- Operational Dashboard: 10 seconds (configurable)

## Key Metrics and Thresholds

### SLA/SLO Thresholds

| Metric | Good | Warning | Critical |
|--------|------|---------|----------|
| Service Success Rate | >99% | 95-99% | <95% |
| CPU Utilization | <70% | 70-90% | >90% |
| Flow Cache Hit Rate | >95% | 90-95% | <90% |
| EMC Hit Rate | >90% | 85-90% | <85% |
| Drop Rate | <0.1% | 0.1-1% | >1% |
| Interface Utilization | <50% | 50-80% | >80% |

### Important Queries

**Tenant Packet Rate**:
```promql
sum by (value) (
  rate(ovs_interface_rx_packets_total[5m]) * on(uuid) group_left(value)
  ovs_interface_external_ids{key="tenant_id"}
)
```

**Per-Tenant Drop Rate**:
```promql
sum by (value) (
  rate(ovs_interface_rx_dropped_total[5m]) * on(uuid) group_left(value)
  ovs_interface_external_ids{key="tenant_id"}
) / sum by (value) (
  rate(ovs_interface_rx_packets_total[5m]) * on(uuid) group_left(value)
  ovs_interface_external_ids{key="tenant_id"}
)
```

**Packet Processing Latency (P95)**:
```promql
histogram_quantile(0.95, sum(rate(ovs_pmd_cycles_per_packet[5m])) by (le))
```

## Troubleshooting

### No Data Appearing

1. **Verify OVS Exporter is running**:
   ```bash
   curl http://localhost:9475/metrics | grep ovs_up
   ```

2. **Check Prometheus targets**:
   - Navigate to Prometheus → Status → Targets
   - Ensure OVS exporter target is "UP"

3. **Verify system_id label**:
   ```promql
   ovs_up
   ```

### Missing Tenant Grouping

1. **Check external_ids configuration**:
   ```bash
   ovs-vsctl list Interface | grep external_ids
   ```

2. **Verify metrics have external_ids labels**:
   ```promql
   ovs_interface_external_ids{key="tenant_id"}
   ```

### High Memory Usage in Grafana

- Reduce query time ranges
- Increase dashboard refresh intervals
- Consider using recording rules in Prometheus for complex queries

## Customization

### Adding New Tenant Identifiers

To use different external_id keys for tenant identification:

1. Update queries to use your key:
   ```promql
   ovs_interface_external_ids{key="your_tenant_key"}
   ```

2. Modify panel titles and descriptions accordingly

### Creating Alerts

Example alert for high drop rate:

```yaml
groups:
  - name: ovs_alerts
    rules:
      - alert: HighDropRate
        expr: |
          sum by (value) (
            rate(ovs_interface_rx_dropped_total[5m]) * on(uuid) group_left(value)
            ovs_interface_external_ids{key="tenant_id"}
          ) > 0.01
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High drop rate for tenant {{ $labels.value }}"
          description: "Tenant {{ $labels.value }} experiencing {{ $value | humanizePercentage }} drop rate"
```

## Best Practices

1. **Regular Updates**: Keep dashboards updated with new OVS exporter metrics
2. **Recording Rules**: Use Prometheus recording rules for complex, frequently-used queries
3. **Access Control**: Use Grafana's folder permissions to control dashboard access
4. **Version Control**: Store dashboard JSON files in version control
5. **Documentation**: Document any custom modifications or site-specific configurations

## Support

For issues related to:
- **Dashboards**: Create an issue in this repository
- **OVS Exporter**: See the main project documentation
- **Grafana**: Consult [Grafana documentation](https://grafana.com/docs/)
- **Prometheus**: Refer to [Prometheus documentation](https://prometheus.io/docs/)