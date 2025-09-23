#!/bin/bash
set -e

MYAPP=ovs-exporter
MYAPP_USER=ovs_exporter
MYAPP_GROUP=ovs_exporter
MYAPP_SERVICE=ovs-exporter
MYAPP_BIN=/usr/sbin/ovs-exporter
MYAPP_DESCRIPTION="Prometheus OVS Exporter"
MYAPP_CONF="/etc/sysconfig/${MYAPP_SERVICE}"
SYSLOG_SERVER=""
SYSLOG_PORT="514"
SYSLOG_PROTOCOL="udp"

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --syslog-server)
            SYSLOG_SERVER="$2"
            shift 2
            ;;
        --syslog-port)
            SYSLOG_PORT="$2"
            shift 2
            ;;
        --syslog-protocol)
            SYSLOG_PROTOCOL="$2"
            shift 2
            ;;
        -h|--help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --syslog-server HOST     Remote syslog server hostname or IP address"
            echo "  --syslog-port PORT       Remote syslog server port (default: 514)"
            echo "  --syslog-protocol PROTO  Protocol to use: udp or tcp (default: udp)"
            echo "  -h, --help               Show this help message"
            echo ""
            echo "Example:"
            echo "  sudo $0 --syslog-server 192.168.1.100 --syslog-port 514"
            exit 0
            ;;
        *)
            echo "ERROR: Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Detect system architecture
ARCH=$(uname -m)
case "$ARCH" in
    x86_64)
        BUILD_ARCH="amd64"
        ;;
    aarch64)
        BUILD_ARCH="arm64"
        ;;
    arm64)
        BUILD_ARCH="arm64"
        ;;
    *)
        echo "ERROR: Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Detect OS (default to linux)
BUILD_OS="linux"

echo "Installing ${MYAPP_DESCRIPTION}..."
echo "Detected architecture: ${BUILD_ARCH}"

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo "ERROR: Please run as root (use sudo)"
    exit 1
fi

# Build if binary doesn't exist
if [ ! -f "./ovs_exporter" ]; then
    if [ ! -f "./bin/${BUILD_OS}-${BUILD_ARCH}/ovs_exporter" ]; then
        echo "Building ovs_exporter for ${BUILD_OS}-${BUILD_ARCH}..."
        make BUILD_OS="${BUILD_OS}" BUILD_ARCH="${BUILD_ARCH}"
    fi
    if [ -f "./bin/${BUILD_OS}-${BUILD_ARCH}/ovs_exporter" ]; then
        cp "./bin/${BUILD_OS}-${BUILD_ARCH}/ovs_exporter" ./ovs_exporter
    else
        echo "ERROR: Could not find or build ovs_exporter binary"
        echo "Please run 'make BUILD_OS=${BUILD_OS} BUILD_ARCH=${BUILD_ARCH}' first or ensure ovs_exporter binary exists"
        exit 1
    fi
fi

# Install binary
echo "Installing binary to ${MYAPP_BIN}..."
rm -rf ${MYAPP_BIN}
cp ./ovs_exporter ${MYAPP_BIN}

# Create group if it doesn't exist
if getent group ${MYAPP_GROUP} >/dev/null; then
    echo "INFO: ${MYAPP_GROUP} group already exists"
else
    echo "INFO: Creating ${MYAPP_GROUP} group..."
    groupadd --system ${MYAPP_GROUP}
fi

# Create user if it doesn't exist
if getent passwd ${MYAPP_USER} >/dev/null; then
    echo "INFO: ${MYAPP_USER} user already exists"
else
    echo "INFO: Creating ${MYAPP_USER} user..."
    useradd --system -d /var/lib/${MYAPP} -s /bin/bash -g ${MYAPP_GROUP} ${MYAPP_USER}
fi

# Add user to openvswitch group
if getent group openvswitch >/dev/null; then
    echo "INFO: Adding ${MYAPP_USER} to openvswitch group..."
    usermod -a -G openvswitch ${MYAPP_USER}
else
    echo "WARNING: openvswitch group does not exist. You may need to add ${MYAPP_USER} to it manually after installing OVS"
fi

# Create working directory
mkdir -p /var/lib/${MYAPP}
chown -R ${MYAPP_USER}:${MYAPP_GROUP} /var/lib/${MYAPP}

# Create systemd service file
echo "Creating systemd service file..."
cat << EOF > /usr/lib/systemd/system/${MYAPP_SERVICE}.service
[Unit]
Description=$MYAPP_DESCRIPTION
After=network.target

[Service]
User=${MYAPP_USER}
Group=${MYAPP_GROUP}
EnvironmentFile=-${MYAPP_CONF}
ExecStart=${MYAPP_BIN} \$OPTIONS
Restart=on-failure

[Install]
WantedBy=multi-user.target
EOF

# Create default configuration file
echo "Creating default configuration file..."
mkdir -p /etc/sysconfig
cat << EOF > ${MYAPP_CONF}
# Configuration for ${MYAPP_DESCRIPTION}
# Add command line options here
OPTIONS="--web.listen-address=:9475 --log.level=info"
EOF

# Set capabilities
echo "Setting capabilities on ${MYAPP_BIN}..."
setcap cap_sys_admin,cap_sys_nice,cap_dac_override+ep ${MYAPP_BIN} || true

# Adjust OVS socket permissions if it exists
if [ -S /var/run/openvswitch/db.sock ]; then
    echo "Adjusting OVS socket permissions..."
    chmod g+w /var/run/openvswitch/db.sock || true
fi

# Configure remote syslog if requested
if [ -n "$SYSLOG_SERVER" ]; then
    echo "Configuring remote syslog forwarding to ${SYSLOG_SERVER}:${SYSLOG_PORT} (${SYSLOG_PROTOCOL})..."

    # Check if rsyslog is installed
    if ! command -v rsyslogd &> /dev/null; then
        echo "WARNING: rsyslog is not installed. Skipping syslog configuration."
        echo "Install rsyslog and re-run the installer with --syslog-server option to enable remote logging."
    else
        # Create rsyslog configuration for ovs-exporter
        RSYSLOG_CONF="/etc/rsyslog.d/30-ovs-exporter.conf"

        # Determine the protocol prefix
        if [ "$SYSLOG_PROTOCOL" = "tcp" ]; then
            PROTO_PREFIX="@@"  # TCP uses @@
        else
            PROTO_PREFIX="@"   # UDP uses @
        fi

        cat << EOF > ${RSYSLOG_CONF}
# Remote syslog forwarding for ${MYAPP_DESCRIPTION}
# Forward all logs from ovs-exporter to remote syslog server

# Match logs from ovs-exporter service
if \$programname == 'ovs-exporter' then {
    # Forward to remote syslog server
    action(type="omfwd"
           target="${SYSLOG_SERVER}"
           port="${SYSLOG_PORT}"
           protocol="${SYSLOG_PROTOCOL}"
           template="RSYSLOG_SyslogProtocol23Format")

    # Optionally stop processing (uncomment to prevent local logging)
    # stop
}

# Alternative: Forward based on systemd unit
if \$programname == 'systemd' and \$msg contains 'ovs-exporter' then {
    ${PROTO_PREFIX}${SYSLOG_SERVER}:${SYSLOG_PORT}
}
EOF

        echo "Created rsyslog configuration at ${RSYSLOG_CONF}"

        # Test rsyslog configuration
        if rsyslogd -N1 -f ${RSYSLOG_CONF} &> /dev/null; then
            echo "Rsyslog configuration is valid"

            # Restart rsyslog to apply changes
            echo "Restarting rsyslog service..."
            systemctl restart rsyslog || service rsyslog restart || true

            echo "Remote syslog forwarding configured successfully"
        else
            echo "WARNING: Rsyslog configuration validation failed. Please check ${RSYSLOG_CONF}"
        fi
    fi
fi

# Reload systemd and enable service
echo "Configuring systemd service..."
systemctl daemon-reload

# Stop service if it's running
systemctl is-active --quiet ${MYAPP_SERVICE} && systemctl stop ${MYAPP_SERVICE}

# Enable and start service
systemctl enable ${MYAPP_SERVICE}
systemctl start ${MYAPP_SERVICE}

# Check if service is running
if systemctl is-active --quiet ${MYAPP_SERVICE}; then
    echo ""
    echo "SUCCESS: ${MYAPP_SERVICE} service is installed and running"
    echo ""
    echo "You can check the status with:"
    echo "  systemctl status ${MYAPP_SERVICE}"
    echo ""
    echo "View logs with:"
    echo "  journalctl -u ${MYAPP_SERVICE} -f"
    echo ""
    if [ -n "$SYSLOG_SERVER" ]; then
        echo "Remote syslog forwarding:"
        echo "  Logs are being forwarded to ${SYSLOG_SERVER}:${SYSLOG_PORT} (${SYSLOG_PROTOCOL})"
        echo "  Configuration: /etc/rsyslog.d/30-ovs-exporter.conf"
        echo ""
    fi
    echo "Metrics are available at:"
    echo "  http://localhost:9475/metrics"
else
    echo ""
    echo "ERROR: ${MYAPP_SERVICE} service is not running"
    echo "Check the status with:"
    echo "  systemctl status ${MYAPP_SERVICE}"
    exit 1
fi