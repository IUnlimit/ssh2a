#!/bin/bash
set -e

OUTPUT_DIR="output"
mkdir -p "$OUTPUT_DIR"

echo "==> Building frontend..."
cd web && pnpm install && pnpm run build && cd ..

echo "==> Compiling binary (linux/amd64)..."
GOOS=linux GOARCH=amd64 go build -o "$OUTPUT_DIR/ssh2a_linux"

echo "==> Build complete: $OUTPUT_DIR/ssh2a_linux"
