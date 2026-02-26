#!/bin/bash
# Self-extracting installer
# Extracts embedded archive and runs install.sh

set -e

echo "================================================"
echo "  myapp Self-extracting Installer"
echo "================================================"
echo ""

INSTALL_DIR="${1:-/opt/myapp}"
TMPDIR=$(mktemp -d)

# Extract embedded tar.gz archive
ARCHIVE_START=$(awk '/^__ARCHIVE_START__$/{print NR + 1; exit 0;}' "$0")
tail -n +"${ARCHIVE_START}" "$0" | tar xz -C "${TMPDIR}"

# Run install script
cd "${TMPDIR}"
chmod +x install.sh
bash install.sh "${INSTALL_DIR}"

# Cleanup
rm -rf "${TMPDIR}"
exit 0

__ARCHIVE_START__
