#!/usr/bin/env bash
set -euo pipefail

echo "INFO: running postinstall"
setcap 'cap_net_bind_service=+ep' /usr/bin/gosmtp
systemctl daemon-reload
