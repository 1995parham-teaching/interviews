#!/usr/bin/env bash
# Starts three backends and three Relay nodes on localhost. Ctrl-C stops all.
#
#   backends  127.0.0.1:9001  127.0.0.1:9002  127.0.0.1:9003
#   nodes     relay-1 :8001   relay-2 :8002   relay-3 :8003
#
# Logs go to ./logs/<name>.log.
set -euo pipefail
cd "$(dirname "$0")/.."

mkdir -p bin logs
go build -o bin/ ./cmd/relay ./cmd/backend

pids=()
trap 'kill "${pids[@]}" 2>/dev/null; wait' EXIT INT TERM

for port in 9001 9002 9003; do
	bin/backend -listen "127.0.0.1:$port" >"logs/backend-$port.log" 2>&1 &
	pids+=($!)
done

backends=127.0.0.1:9001,127.0.0.1:9002,127.0.0.1:9003
nodes=(127.0.0.1:8001 127.0.0.1:8002 127.0.0.1:8003)
for i in 0 1 2; do
	peers=$(printf '%s\n' "${nodes[@]}" | grep -vx "${nodes[$i]}" | paste -sd, -)
	bin/relay -id "relay-$((i + 1))" -listen "${nodes[$i]}" -peers "$peers" -backends "$backends" \
		>"logs/relay-$((i + 1)).log" 2>&1 &
	pids+=($!)
done

echo "cluster up: nodes on :8001 :8002 :8003, backends on :9001 :9002 :9003; logs in ./logs"
wait
