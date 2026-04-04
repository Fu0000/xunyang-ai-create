#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
LOG_DIR="$ROOT_DIR/logs"
PID_DIR="$ROOT_DIR/.pids"

mkdir -p "$LOG_DIR" "$PID_DIR"

if [[ -f "$PID_DIR/frontend-admin.pid" ]] && kill -0 "$(cat "$PID_DIR/frontend-admin.pid")" 2>/dev/null; then
  echo "frontend-admin already running: PID $(cat "$PID_DIR/frontend-admin.pid")"
  exit 0
fi

cd "$ROOT_DIR/frontend-admin"
nohup npm run dev -- --host 127.0.0.1 --port 5174 > "$LOG_DIR/frontend-admin.log" 2>&1 &
ADMIN_PID=$!
echo "$ADMIN_PID" > "$PID_DIR/frontend-admin.pid"

echo "frontend-admin started: PID $ADMIN_PID"
echo "log: $LOG_DIR/frontend-admin.log"
