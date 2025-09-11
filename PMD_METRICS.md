# PMD Performance Metrics

This exporter now supports comprehensive PMD (Poll Mode Driver) performance metrics for DPDK-enabled Open vSwitch deployments, as described in the [Red Hat blog post on OVS observability features](https://www.redhat.com/en/blog/amazing-new-observability-features-open-vswitch).

## Supported PMD Metrics

The exporter automatically collects comprehensive PMD performance metrics when running on DPDK-enabled OVS deployments. These metrics are critical for understanding and optimizing high-performance packet processing.

### Core PMD Performance Metrics

| Metric | Description | Labels |
|--------|-------------|--------|
| `ovs_pmd_cpu_utilization_percent` | CPU utilization percentage of PMD thread | `system_id`, `pmd_id`, `numa_id`, `core_id` |
| `ovs_pmd_cycles_per_iteration` | Average cycles spent per PMD iteration | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_packets_per_iteration` | Average packets processed per PMD iteration | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_cycles_per_packet` | Average cycles spent per packet in PMD | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_packets_per_batch` | Average packets per batch in PMD | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_max_vhost_queue_length` | Maximum vhost queue length observed | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_avg_vhost_queue_length` | Average vhost queue length | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_vhost_queue_full_total` | Number of times vhost queue was full | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_iterations_total` | Total number of PMD iterations | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_sleep_iterations_total` | Total sleep iterations for PMD thread | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_busy_cycles_total` | Total cycles where PMD was busy | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_idle_cycles_total` | Total idle cycles for PMD thread | `system_id`, `pmd_id`, `numa_id` |

### RX/TX Batch Statistics

| Metric | Description | Labels |
|--------|-------------|--------|
| `ovs_pmd_rx_batches_total` | Total number of RX batches processed | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_rx_packets_total` | Total number of RX packets processed | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_avg_rx_batch_size` | Average RX batch size | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_max_rx_batch_size` | Maximum RX batch size observed | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_tx_batches_total` | Total number of TX batches processed | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_tx_packets_total` | Total number of TX packets processed | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_avg_tx_batch_size` | Average TX batch size | `system_id`, `pmd_id`, `numa_id` |

### Flow Lookup Statistics

| Metric | Description | Labels |
|--------|-------------|--------|
| `ovs_pmd_exact_match_hit_total` | Number of exact match hits | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_masked_hit_total` | Number of masked hits | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_miss_total` | Number of flow misses | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_lost_total` | Number of lost packets | `system_id`, `pmd_id`, `numa_id` |

### Suspicious Iteration Detection

| Metric | Description | Labels |
|--------|-------------|--------|
| `ovs_pmd_suspicious_iterations_total` | Number of suspicious iterations detected | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_suspicious_iterations_percent` | Percentage of iterations that are suspicious | `system_id`, `pmd_id`, `numa_id` |

### Upcall Metrics

| Metric | Description | Labels |
|--------|-------------|--------|
| `ovs_pmd_upcalls_total` | Total number of upcalls from PMD | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_upcall_cycles_total` | Total cycles spent in upcalls | `system_id`, `pmd_id`, `numa_id` |

### vHost Specific Counters

| Metric | Description | Labels |
|--------|-------------|--------|
| `ovs_vhost_tx_retries_total` | Total number of vhost transmit retries | `system_id`, `pmd_id`, `numa_id` |
| `ovs_vhost_tx_contention_total` | Total number of vhost transmit contentions | `system_id`, `pmd_id`, `numa_id` |
| `ovs_vhost_tx_irqs_total` | Total number of vhost transmit IRQs | `system_id`, `pmd_id`, `numa_id` |

## Packet Drop Statistics

The exporter collects comprehensive packet drop statistics through two mechanisms:

### Individual Drop Counter Metrics

The exporter now exposes all 24 specific drop counters as individual metrics via `ovs_datapath_drops_total` with a `drop_reason` label:

| Drop Reason | Description |
|-------------|-------------|
| `datapath_drop_upcall_error` | Upcall error drops |
| `datapath_drop_lock_error` | Lock contention drops |
| `datapath_drop_rx_invalid_packet` | Invalid RX packet drops |
| `datapath_drop_meter` | Meter-based drops |
| `datapath_drop_userspace_action_error` | Userspace action errors |
| `datapath_drop_tunnel_push_error` | Tunnel push errors |
| `datapath_drop_tunnel_pop_error` | Tunnel pop errors |
| `datapath_drop_recirc_error` | Recirculation errors |
| `datapath_drop_invalid_port` | Invalid port drops |
| `datapath_drop_invalid_tnl_port` | Invalid tunnel port drops |
| `datapath_drop_sample_error` | Sampling errors |
| `datapath_drop_nsh_decap_error` | NSH decapsulation errors |
| `drop_action_of_pipeline` | OpenFlow pipeline drops |
| `drop_action_bridge_not_found` | Bridge not found drops |
| `drop_action_recursion_too_deep` | Recursion limit exceeded |
| `drop_action_too_many_resubmit` | Too many resubmits |
| `drop_action_stack_too_deep` | Stack overflow |
| `drop_action_no_recirculation` | No recirculation available |
| `drop_action_recirculation_conflict` | Recirculation conflicts |
| `drop_action_too_many_mpls_labels` | MPLS label limit exceeded |
| `drop_action_invalid_tunnel_metadata` | Invalid tunnel metadata |
| `drop_action_unsupported_packet_type` | Unsupported packet type |
| `drop_action_congestion` | Congestion drops |
| `drop_action_forwarding_disabled` | Forwarding disabled |

### Coverage-based Drop Statistics

Additionally, all drop statistics are available through the existing coverage metrics mechanism as `ovs_coverage_total` metrics with event labels.

## Example Queries

### Monitor PMD Efficiency
```promql
# Cycles per packet by PMD thread
ovs_pmd_cycles_per_packet{job="ovs-exporter"}

# Packet processing rate
rate(ovs_pmd_iterations_total[5m]) * ovs_pmd_packets_per_iteration
```

### Identify Performance Issues
```promql
# High upcall rate (potential flow miss issues)
rate(ovs_pmd_upcalls_total[5m]) > 1000

# vHost contention issues
rate(ovs_vhost_tx_contention_total[5m]) > 0
```

### Monitor Packet Drops
```promql
# Total packet drops by type
increase(ovs_coverage_total{event=~".*drop.*"}[5m])

# Specific drop reasons
ovs_coverage_total{event="datapath_drop_upcall_error"}
```

## Data Collection

PMD metrics are collected using the following ovs-appctl commands:
- `ovs-appctl dpif-netdev/pmd-perf-show` - Detailed PMD performance statistics
- `ovs-appctl dpif-netdev/pmd-stats-show` - Additional PMD statistics
- `ovs-appctl coverage/show` - Coverage counters including drop statistics

These commands are automatically executed during each collection interval. On non-DPDK deployments, PMD metric collection is gracefully skipped.

## Configuration

No additional configuration is required. The exporter automatically detects DPDK deployments and collects PMD metrics when available. The standard polling interval applies to PMD metrics collection.

## Grafana Dashboard

PMD metrics can be visualized using Grafana dashboards. Key panels to consider:

1. **PMD Performance Overview**
   - Cycles per packet (line graph)
   - Packets per iteration (gauge)
   - PMD utilization (percentage)

2. **vHost Statistics**
   - TX retries over time
   - Contention events
   - IRQ distribution

3. **Packet Drop Analysis**
   - Drop reasons breakdown (pie chart)
   - Drop rate trends (line graph)
   - Top drop causes (table)

## Troubleshooting

### No PMD Metrics Available

If PMD metrics are not appearing:

1. Verify DPDK is enabled in OVS:
   ```bash
   ovs-vsctl get Open_vSwitch . other_config:dpdk-init
   ```

2. Check if PMD threads are running:
   ```bash
   ovs-appctl dpif-netdev/pmd-stats-show
   ```

3. Ensure the exporter has sufficient permissions to run ovs-appctl commands.

### High Drop Rates

Use coverage metrics to identify drop causes:
```bash
ovs-appctl coverage/show | grep drop
```

Common causes:
- Flow table misses (check flow rules)
- MTU issues (check interface MTU settings)
- QoS limits (review QoS configuration)