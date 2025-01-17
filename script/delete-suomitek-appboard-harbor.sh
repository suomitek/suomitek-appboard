#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

# Constants
ROOT_DIR="$(cd "$( dirname "${BASH_SOURCE[0]}" )/.." >/dev/null && pwd)"

# Load Libraries
# shellcheck disable=SC1090
. "${ROOT_DIR}/script/libtest.sh"
# shellcheck disable=SC1090
. "${ROOT_DIR}/script/liblog.sh"

# Delete Harbor
info "---------------------"
info "-- Harbor deletion --"
info "---------------------"
echo
"$ROOT_DIR"/script/delete-harbor.sh --namespace "harbor"
# Delete Suomitek-appboard
info "-----------------------"
info "-- Suomitek-appboard deletion --"
info "-----------------------"
echo
"$ROOT_DIR"/script/delete-suomitek-appboard.sh --namespace "suomitek-appboard"
