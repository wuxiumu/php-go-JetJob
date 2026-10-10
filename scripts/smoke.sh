#!/usr/bin/env bash
# 不连接 MySQL / Redis。只跑不依赖数据库的 Go 测试。
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT/scheduler-go1.24.2"
go test ./internal/model ./internal/utils -count=1
go build -o /dev/null ./cmd/scheduler
echo "smoke ok"
