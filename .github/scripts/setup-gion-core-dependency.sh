#!/usr/bin/env bash
set -euo pipefail

if [[ -d "../gion-core" ]]; then
  exit 0
fi

repo="${GION_CORE_REPO:-https://github.com/tasuku43/gion-core.git}"
ref="${GION_CORE_REF:-main}"

git clone --depth=1 --branch "${ref}" "${repo}" ../gion-core
