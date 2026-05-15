#!/usr/bin/env bash
set -euo pipefail

if ! command -v golangci-lint >/dev/null 2>&1; then
  echo "Instalando golangci-lint..."
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0
fi

golangci-lint run ./...
