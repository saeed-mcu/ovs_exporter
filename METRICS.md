# OVS Exporter Metrics

This document provides a comprehensive reference for all metrics exported by the OVS Exporter. All metrics follow [Prometheus naming best practices](https://prometheus.io/docs/practices/naming/).

## Table of Contents

- [System Metrics](#system-metrics)
- [Process and Component Metrics](#process-and-component-metrics)
- [Coverage and Memory Metrics](#coverage-and-memory-metrics)
- [Datapath Metrics](#datapath-metrics)
- [Interface Metrics](#interface-metrics)
- [PMD Performance Metrics](#pmd-performance-metrics)
- [Flow Cache Metrics](#flow-cache-metrics)
- [vHost Metrics](#vhost-metrics)
- [Drop Statistics](#drop-statistics)

## System Metrics

### Core System Status

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `ovs_up` | Gauge | Is OVN stack up (1) or is it down (0) | - |
| `ovs_info` | Gauge | Basic information about OVN stack. Always set to 1 | `system_id`, `rundir`, `hostname`, `system_type`, `system_version`, `ovs_version`, `db_version` |
| `ovs_failed_requests_total` | Counter | The number of failed requests to OVN stack | `system_id` |
| `ovs_next_poll_timestamp_seconds` | Gauge | The timestamp of the next potential poll of OVN stack | `system_id` |
| `ovs_exporter_build_info` | Gauge | Build information about the exporter itself | `version`, `revision`, `branch`, `goversion` |

## Process and Component Metrics

### Process Information

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `ovs_pid` | Gauge | The process ID of a running OVN component (0 if not running) | `system_id`, `component`, `user`, `group` |
| `ovs_network_port_up` | Gauge | Whether the network port is up (1) or down (0) for database connection | `system_id`, `component`, `usage` |

### Log Files

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `ovs_log_file_size_bytes` | Gauge | The size of a log file associated with an OVN component | `system_id`, `component`, `filename` |
| `ovs_log_events` | Gauge | The number of recorded log messages by severity and source | `system_id`, `component`, `severity`, `source` |

### Database Files

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `ovs_db_file_size_bytes` | Gauge | The size of a database file associated with an OVN component | `system_id`, `component`, `filename` |

## Coverage and Memory Metrics

### Coverage Statistics

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `ovs_coverage_total` | Counter | Total number of times particular events occur during OVSDB daemon runtime | `system_id`, `component`, `event` |
| `ovs_coverage_avg` | Gauge | Average rate of events occurring during OVSDB daemon runtime | `system_id`, `component`, `event`, `interval` |

### Memory Usage

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `ovs_memory_usage_bytes` | Gauge | Memory usage in bytes | `system_id`, `component`, `facility` |

## Datapath Metrics

### Datapath Configuration

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `ovs_dp_interface` | Gauge | Represents an existing datapath interface (always 1) | `system_id`, `datapath`, `bridge`, `name`, `ofport`, `index`, `port_type` |
| `ovs_dp_bridge_interfaces` | Gauge | The number of interfaces attached to a bridge | `system_id`, `datapath`, `bridge` |
| `ovs_dp_flows` | Gauge | The number of flows in a datapath | `system_id`, `datapath` |

### Datapath Lookups

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `ovs_dp_lookups_hit_total` | Counter | Packets matching existing flows in the datapath | `system_id`, `datapath` |
| `ovs_dp_lookups_missed_total` | Counter | Packets not matching any existing flow | `system_id`, `datapath` |
| `ovs_dp_lookups_lost_total` | Counter | Packets destined for userspace but dropped before reaching it | `system_id`, `datapath` |

### Datapath Masks

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `ovs_dp_masks_hit_total` | Counter | Total number of masks visited for matching incoming packets | `system_id`, `datapath` |
| `ovs_dp_masks_total` | Counter | The number of masks in a datapath | `system_id`, `datapath` |
| `ovs_dp_masks_hit_ratio` | Gauge | Average number of masks visited per packet | `system_id`, `datapath` |

## Interface Metrics

### Interface Status

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `ovs_interface` | Gauge | Primary metric for OVS interface (always 1) | `system_id`, `uuid`, `name` |
| `ovs_interface_admin_state` | Gauge | Administrative state: down(0), up(1), other(2) | `system_id`, `uuid` |
| `ovs_interface_link_state` | Gauge | Observed link state: down(0), up(1), other(2) | `system_id`, `uuid` |
| `ovs_interface_duplex` | Gauge | Duplex mode: other(0), half(1), full(2) | `system_id`, `uuid` |
| `ovs_interface_mac_in_use` | Gauge | MAC address in use (always 1) | `system_id`, `uuid`, `mac_address` |

### Interface Configuration

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `ovs_interface_mtu_bytes` | Gauge | Currently configured MTU in bytes | `system_id`, `uuid` |
| `ovs_interface_ingress_policing_rate_kilobits_per_second` | Gauge | Maximum ingress rate in kbps (0 = disabled) | `system_id`, `uuid` |
| `ovs_interface_ingress_policing_burst_kilobits` | Gauge | Maximum burst size in kb (default 8000 if 0) | `system_id`, `uuid` |
| `ovs_interface_link_speed_bits_per_second` | Gauge | Negotiated link speed in bps | `system_id`, `uuid` |
| `ovs_interface_openflow_port` | Gauge | OpenFlow port ID | `system_id`, `uuid` |
| `ovs_interface_index` | Gauge | Interface index | `system_id`, `uuid` |
| `ovs_interface_local_index` | Gauge | Local index | `system_id`, `uuid` |

### Interface Statistics - Receive

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `ovs_interface_rx_packets_total` | Counter | Number of received packets | `system_id`, `uuid` |
| `ovs_interface_rx_bytes` | Counter | Number of received bytes | `system_id`, `uuid` |
| `ovs_interface_rx_multicast_packets_total` | Counter | Number of received multicast packets | `system_id`, `uuid` |

### Interface Statistics - Receive Errors

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `ovs_interface_rx_dropped_total` | Counter | Number of input packets dropped | `system_id`, `uuid` |
| `ovs_interface_rx_errors_total` | Counter | Total number of receive errors | `system_id`, `uuid` |
| `ovs_interface_rx_crc_errors_total` | Counter | Number of CRC errors | `system_id`, `uuid` |
| `ovs_interface_rx_frame_errors_total` | Counter | Number of frame alignment errors | `system_id`, `uuid` |
| `ovs_interface_rx_overrun_errors_total` | Counter | Number of RX overrun errors | `system_id`, `uuid` |
| `ovs_interface_rx_missed_errors_total` | Counter | Number of missed packets | `system_id`, `uuid` |

### Interface Statistics - Transmit

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `ovs_interface_tx_packets_total` | Counter | Number of transmitted packets | `system_id`, `uuid` |
| `ovs_interface_tx_bytes` | Counter | Number of transmitted bytes | `system_id`, `uuid` |
| `ovs_interface_tx_dropped_total` | Counter | Number of output packets dropped | `system_id`, `uuid` |
| `ovs_interface_tx_errors_total` | Counter | Total number of transmit errors | `system_id`, `uuid` |
| `ovs_interface_collisions_total` | Counter | Number of collisions | `system_id`, `uuid` |

### Interface Link Events

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `ovs_interface_link_resets_total` | Counter | Number of times link state changed | `system_id`, `uuid` |

### Interface Key-Value Pairs

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `ovs_interface_status` | Gauge | Port status key-value pairs (always 1) | `system_id`, `uuid`, `key`, `value` |
| `ovs_interface_options` | Gauge | Interface options key-value pairs (always 1) | `system_id`, `uuid`, `key`, `value` |
| `ovs_interface_external_ids` | Gauge | External IDs key-value pairs (always 1) | `system_id`, `uuid`, `key`, `value` |

## PMD Performance Metrics

PMD (Poll Mode Driver) metrics are available for DPDK-enabled OVS deployments.

### Core PMD Performance

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `ovs_pmd_cpu_utilization_ratio` | Gauge | CPU utilization ratio of PMD thread (0-1) | `system_id`, `pmd_id`, `numa_id`, `core_id` |
| `ovs_pmd_cycles_per_iteration` | Gauge | Average cycles spent per PMD iteration | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_packets_per_iteration` | Gauge | Average packets processed per iteration | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_cycles_per_packet` | Gauge | Average cycles spent per packet | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_packets_per_batch` | Gauge | Average packets per batch | `system_id`, `pmd_id`, `numa_id` |

### PMD Iteration Statistics

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `ovs_pmd_iterations_total` | Counter | Total number of PMD iterations | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_sleep_iterations_total` | Counter | Total sleep iterations | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_busy_cycles_total` | Counter | Total cycles where PMD was busy | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_idle_cycles_total` | Counter | Total idle cycles | `system_id`, `pmd_id`, `numa_id` |

### RX/TX Batch Statistics

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `ovs_pmd_rx_batches_total` | Counter | Total number of RX batches processed | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_rx_packets_total` | Counter | Total number of RX packets processed | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_avg_rx_batch_size` | Gauge | Average RX batch size | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_max_rx_batch_size` | Gauge | Maximum RX batch size observed | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_tx_batches_total` | Counter | Total number of TX batches processed | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_tx_packets_total` | Counter | Total number of TX packets processed | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_avg_tx_batch_size` | Gauge | Average TX batch size | `system_id`, `pmd_id`, `numa_id` |

### Flow Lookup Statistics

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `ovs_pmd_exact_match_hit_total` | Counter | Number of exact match hits | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_masked_hit_total` | Counter | Number of masked hits | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_miss_total` | Counter | Number of flow misses | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_lost_total` | Counter | Number of lost packets | `system_id`, `pmd_id`, `numa_id` |

### Suspicious Iterations

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `ovs_pmd_suspicious_iterations_total` | Counter | Number of suspicious iterations detected | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_suspicious_iterations_ratio` | Gauge | Ratio of suspicious iterations (0-1) | `system_id`, `pmd_id`, `numa_id` |

### Upcall Metrics

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `ovs_pmd_upcalls_total` | Counter | Total number of upcalls from PMD | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_upcall_cycles_total` | Counter | Total cycles spent in upcalls | `system_id`, `pmd_id`, `numa_id` |

## Flow Cache Metrics

### Cache Performance

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `ovs_flow_cache_emc_hit_ratio` | Gauge | Exact Match Cache hit ratio (0-1) | `system_id`, `pmd_id`, `numa_id` |
| `ovs_flow_cache_emc_hits_total` | Counter | Total EMC hits | `system_id`, `pmd_id`, `numa_id` |
| `ovs_flow_cache_emc_inserts_total` | Counter | Total EMC insertions | `system_id`, `pmd_id`, `numa_id` |
| `ovs_flow_cache_smc_hit_ratio` | Gauge | Signature Match Cache hit ratio (0-1) | `system_id`, `pmd_id`, `numa_id` |
| `ovs_flow_cache_smc_hits_total` | Counter | Total SMC hits | `system_id`, `pmd_id`, `numa_id` |
| `ovs_flow_cache_megaflow_hit_ratio` | Gauge | Megaflow cache hit ratio (0-1) | `system_id`, `pmd_id`, `numa_id` |
| `ovs_flow_cache_megaflow_hits_total` | Counter | Total Megaflow cache hits | `system_id`, `pmd_id`, `numa_id` |
| `ovs_flow_cache_megaflow_misses_total` | Counter | Total Megaflow cache misses | `system_id`, `pmd_id`, `numa_id` |
| `ovs_flow_cache_lookups_total` | Counter | Total flow cache lookups | `system_id`, `pmd_id`, `numa_id` |

## vHost Metrics

### vHost Queue Statistics

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `ovs_pmd_max_vhost_queue_length` | Gauge | Maximum vhost queue length observed | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_avg_vhost_queue_length` | Gauge | Average vhost queue length | `system_id`, `pmd_id`, `numa_id` |
| `ovs_pmd_vhost_queue_full_total` | Counter | Number of times vhost queue was full | `system_id`, `pmd_id`, `numa_id` |

### vHost Transmission

| Metric | Type | Description | Labels |
|--------|------|-------------|--------|
| `ovs_vhost_tx_retries_total` | Counter | Total number of vhost transmit retries | `system_id`, `pmd_id`, `numa_id` |
| `ovs_vhost_tx_contention_total` | Counter | Total number of vhost transmit contentions | `system_id`, `pmd_id`, `numa_id` |
| `ovs_vhost_tx_irqs_total` | Counter | Total number of vhost transmit IRQs | `system_id`, `pmd_id`, `numa_id` |

## Drop Statistics

### Datapath Drops

The exporter exposes detailed drop counters via `ovs_datapath_drops_total` with a `drop_reason` label:

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

Note: Drop statistics are also available through coverage metrics (`ovs_coverage_total`) with event labels.

## Example Queries

### System Health
```promql
# Check if OVS is up
ovs_up

# Monitor failed requests
rate(ovs_failed_requests_total[5m])
```

### Interface Performance
```promql
# Interface traffic rate (packets/sec)
rate(ovs_interface_rx_packets_total[5m])

# Interface error rate
rate(ovs_interface_rx_errors_total[5m])

# Interface utilization (requires link speed)
rate(ovs_interface_rx_bytes[5m]) * 8 / ovs_interface_link_speed_bits_per_second
```

### PMD Efficiency
```promql
# PMD CPU utilization
ovs_pmd_cpu_utilization_ratio

# Cycles per packet by PMD thread
ovs_pmd_cycles_per_packet

# Packet processing rate
rate(ovs_pmd_iterations_total[5m]) * ovs_pmd_packets_per_iteration
```

### Flow Cache Performance
```promql
# EMC hit ratio
ovs_flow_cache_emc_hit_ratio

# Cache miss rate
1 - ovs_flow_cache_megaflow_hit_ratio

# Total cache lookups rate
rate(ovs_flow_cache_lookups_total[5m])
```

### Drop Analysis
```promql
# Total drops by reason
increase(ovs_datapath_drops_total[5m])

# Top drop reasons
topk(5, increase(ovs_datapath_drops_total[5m]))

# Drop rate trend
rate(ovs_coverage_total{event=~".*drop.*"}[5m])
```

### vHost Performance
```promql
# vHost TX retry rate
rate(ovs_vhost_tx_retries_total[5m])

# vHost queue saturation
ovs_pmd_avg_vhost_queue_length / ovs_pmd_max_vhost_queue_length

# vHost contention issues
rate(ovs_vhost_tx_contention_total[5m]) > 0
```

## Data Collection

The exporter collects metrics using various methods:

### OVS Commands
- `ovs-appctl dpif/show` - Datapath interfaces
- `ovs-appctl dpif-netdev/pmd-perf-show` - PMD performance statistics
- `ovs-appctl dpif-netdev/pmd-stats-show` - Additional PMD statistics
- `ovs-appctl coverage/show` - Coverage counters including drops
- `ovs-appctl memory/show` - Memory usage statistics

### Database Queries
- Direct queries to Open_vSwitch database via Unix socket
- Interface statistics from Interface table
- System information from Open_vSwitch table

### File System
- Log file sizes from `/var/log/openvswitch/`
- Database file sizes from `/etc/openvswitch/`
- Process information from `/var/run/openvswitch/`

## Configuration

The exporter uses standard configuration options:

- **Polling Interval**: Default 15 seconds (configurable via `-poll.interval`)
- **OVS Socket**: Default `/var/run/openvswitch/db.sock` (configurable)
- **Metrics Port**: Default 9475 (configurable via `-web.listen-address`)

PMD metrics are automatically detected and collected when available. No additional configuration is required for DPDK deployments.

## Troubleshooting

### No Metrics Available
1. Verify OVS is running: `systemctl status openvswitch`
2. Check socket permissions: `ls -l /var/run/openvswitch/db.sock`
3. Verify exporter has correct capabilities: `getcap /usr/sbin/ovs-exporter`

### Missing PMD Metrics
1. Verify DPDK is enabled: `ovs-vsctl get Open_vSwitch . other_config:dpdk-init`
2. Check PMD threads: `ovs-appctl dpif-netdev/pmd-stats-show`
3. Ensure sufficient permissions for ovs-appctl commands

### High Drop Rates
1. Check coverage metrics: `ovs-appctl coverage/show | grep drop`
2. Review flow rules for mismatches
3. Verify MTU settings across interfaces
4. Check QoS configuration

## References

- [Open vSwitch Documentation](http://www.openvswitch.org/support/dist-docs/)
- [Prometheus Naming Best Practices](https://prometheus.io/docs/practices/naming/)
- [Red Hat: OVS Observability Features](https://www.redhat.com/en/blog/amazing-new-observability-features-open-vswitch)
- [OVS Database Schema](http://www.openvswitch.org/support/dist-docs/ovs-vswitchd.conf.db.5.html)