# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Prometheus exporter for Open Virtual Switch (OVS) that exports metrics from OVS vswitchd service, Open_vSwitch database, and OVN ovn-controller service.

**Requirements**: Go 1.24 or later

## Build Commands

```bash
# Build the binary (default: linux-amd64)
make

# Build for specific architecture
make BUILD_OS="linux" BUILD_ARCH="arm64"

# Run quick test (requires sudo)
make qtest

# Run full tests with coverage (requires sudo and OVS installed)
make test
make coverage

# Deploy locally (requires sudo)
make deploy

# Create distribution package
make dist

# Clean build artifacts
make clean
```

## Testing

Tests require Open vSwitch to be installed and running. The test suite uses sudo to interact with OVS:

```bash
# Run tests (requires sudo)
make test

# Generate coverage report
make coverage
```

## Release Process

```bash
# Create a new release (must be on main branch with clean git state)
make release
```

## Architecture

The exporter follows a standard Prometheus exporter pattern:

- **Entry Point**: `cmd/ovs_exporter/main.go` - CLI interface that sets up the HTTP server and configures the exporter
- **Core Exporter**: `pkg/ovs_exporter/ovs_exporter.go` - Implements the Prometheus collector interface, connects to OVS via JSON-RPC unix socket
- **PMD Metrics**: `pkg/ovs_exporter/pmd_metrics.go` - Collects DPDK PMD performance metrics (cycles, packets, vhost stats)
- **Dependencies**: Uses `github.com/greenpau/ovsdb` for OVS database interaction

Key configuration paths (configurable via CLI flags):
- OVS socket: `/var/run/openvswitch/db.sock`
- OVS database: `/etc/openvswitch/conf.db`
- Log files: `/var/log/openvswitch/`
- PID files: `/var/run/openvswitch/`

The exporter polls OVS at configurable intervals (default: 15 seconds) and exposes metrics on port 9475 by default.

## Development Workflow

1. The exporter requires appropriate permissions to access OVS sockets and files
2. Uses `setcap` for capabilities: `cap_sys_admin,cap_sys_nice,cap_dac_override+ep`
3. Systemd service runs as `ovs_exporter` user with membership in `openvswitch` group
4. Binary is installed to `/usr/sbin/ovs-exporter` when deployed