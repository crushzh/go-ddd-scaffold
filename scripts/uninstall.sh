#!/bin/bash
# Uninstall script
# Usage: ./uninstall.sh [install_dir]

set -e

APP_NAME="myapp"
INSTALL_DIR="${1:-/opt/${APP_NAME}}"

echo "================================================"
echo "  ${APP_NAME} Uninstall"
echo "================================================"
echo ""

if [ "$(id -u)" -ne 0 ]; then
    echo "[ERROR] Please run as root"
    exit 1
fi

# Stop and disable service
echo "[1/3] Stopping service..."
systemctl stop ${APP_NAME}.service 2>/dev/null || true
systemctl disable ${APP_NAME}.service 2>/dev/null || true

# Remove service file
echo "[2/3] Removing system service..."
rm -f /etc/systemd/system/${APP_NAME}.service
systemctl daemon-reload

# Remove program files (preserve data directory)
echo "[3/3] Removing program files..."
rm -f "${INSTALL_DIR}/${APP_NAME}"
rm -rf "${INSTALL_DIR}/configs"
rm -rf "${INSTALL_DIR}/scripts"

echo ""
echo "Uninstall complete!"
echo "Note: Data directory ${INSTALL_DIR}/data has been preserved."
echo "      Remove it manually if no longer needed."
echo ""
