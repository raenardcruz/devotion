#!/bin/bash
# Generate Traditional Latin Mass (Divinum Officium) PDF with Orthodox Prayer Book styling
# Usage:
#   ./generate_missa.sh              (defaults to today)
#   ./generate_missa.sh 09-21-2026   (MM-DD-YYYY)
#   ./generate_missa.sh 2026-12-25   (YYYY-MM-DD)

set -e
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
node "$DIR/scripts/generate_missa_pdf.mjs" "$@"
