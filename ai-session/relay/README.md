# Relay

Relay is a distributed HTTP load balancer. Several Relay nodes run side by side, each one able to take any client request. They agree on the list of backends by gossip and route with a consistent-hash ring, so every node sends the same client to the same backend, and a session stays on the backend that holds its state.

Only the Go standard library is used. Go 1.27 or later.

## How it works

- **Ring.** Every active backend is placed on a consistent-hash ring, 64 points per unit of weight. A request goes to the backend that owns its affinity key: the `X-Relay-Key` header if the client sends one, otherwise the client's IP address. Behind a trusted proxy (`-trusted-proxies`), the client address comes from `X-Forwarded-For`.
- **Session pins.** When a backend joins, it takes over part of the ring. To keep existing sessions where their state is, a node pins each key to the backend that served it and keeps it there while that backend is usable. Pins expire when a key is idle for `-pin-ttl`.
- **Health.** Each node checks every backend with `GET /health`, once a second. Three failed checks in a row mark a backend down; two successful checks in a row bring it back. Health is not shared between nodes.
- **Retries.** A request that fails is tried once more on another backend when that is safe: always for GET, HEAD and OPTIONS, and for other methods only when the first backend never received the request.
- **Backend list.** All nodes start with the same seed list. An operator can add, drain, reweight or remove a backend through any node. Every change carries a version from the node's logical clock, and gossip spreads it to the other nodes within a few seconds. A removed backend is remembered as removed for `-tombstone-ttl`.

```
            ┌──────────┐   gossip   ┌──────────┐   gossip   ┌──────────┐
 clients ──▶│ relay-1  │◀──────────▶│ relay-2  │◀──────────▶│ relay-3  │
            └────┬─────┘            └────┬─────┘            └────┬─────┘
                 │   ring, pins, health  │                       │
                 ▼                       ▼                       ▼
            backend :9001          backend :9002           backend :9003
```

## Layout

```
cmd/relay      one Relay node
cmd/backend    a demo backend to put behind it
cmd/cluster    runs three backends and three nodes on this machine
cmd/relayctl   pokes a running cluster
config         defaults, the JSON file and flags
ring           the consistent-hash ring
affinity       session pins
health         active health checks
cluster        the shared backend list, its versions, and gossip
proxy          forwarding, client addresses and retries
metrics        counters in the Prometheus text format
node           wires the above together; the admin API
```

## Running it

One terminal runs the cluster:

```
go run ./cmd/cluster
```

It builds `bin/relay`, `bin/backend` and `bin/relayctl`, starts three backends on ports 9001 to 9003 and three nodes, `relay-1` to `relay-3`, on ports 8001 to 8003, and writes their logs to `logs/`. For the demo it shortens the tombstone TTL to 30 seconds and trusts `X-Forwarded-For` from 127.0.0.1. It then takes commands:

```
ps                 what is running
stop relay-3       stop one process; start relay-3 brings it back
restart 9002       names, ports and relay2-style shorthands all work
logs relay-1       the last lines of its log
quit               stop everything (Ctrl-C too)
```

A second terminal talks to it:

```
bin/relayctl status                         every backend as each node sees it, and each node's runtime
bin/relayctl send -n 20                     20 requests with keys user-0 .. user-19, and where they went
bin/relayctl send -n 20 -node all           the same through every node
bin/relayctl send -n 6 -key ''              without X-Relay-Key
bin/relayctl send -method POST -body '{"cart":1}' /orders
bin/relayctl drain relay-2 9002             also: activate, weight, remove
bin/relayctl fail 9002 health               also: fail 9002 rate 0.5 | crash 0.5 | delay 300ms
bin/relayctl heal 9002
bin/relayctl orders                         also: orders reset
bin/relayctl help
```

`go run ./cmd/relayctl ...` works the same, and so does plain `curl`:

```
curl -H 'X-Relay-Key: alice' localhost:8001/hello
curl localhost:8001/relay/status
```

If the ports are taken: `go run ./cmd/cluster -node-port 18001 -backend-port 19001`; `relayctl` finds the new ports in `logs/cluster.json`. A single node by hand:

```
go run ./cmd/relay -id relay-1 -listen 127.0.0.1:8001 \
    -peers 127.0.0.1:8002,127.0.0.1:8003 \
    -backends 127.0.0.1:9001,127.0.0.1:9002=2,127.0.0.1:9003
```

`go run ./cmd/relay -h` lists every flag; `-config relay.json` reads the same settings from a file (see `config/config.go` for the field names).

## Node API

Everything under `/relay/` is the node's own API; every other path is forwarded to a backend.

| Request | Effect |
| --- | --- |
| `GET /relay/status` | The node, its clock, ring, peers, pins and runtime, and every known backend with its status, weight, version and health as this node sees it |
| `PUT /relay/backends/{ip:port}` with `{"status":"active"\|"draining","weight":1..10}` | Adds the backend or changes it; a field left out keeps its value |
| `DELETE /relay/backends/{ip:port}` | Removes the backend |
| `GET /relay/metrics` | Counters in the Prometheus text format |
| `GET /relay/debug/pprof/` | The Go profiler |
| `POST /relay/gossip` | Used between nodes |

With `-admin-token`, `PUT` and `DELETE` need `Authorization: Bearer <token>`. Responses that went through Relay carry `X-Relay-Node`, `X-Relay-Backend` and `X-Request-Id`.

## Demo backend

`cmd/backend` answers every request with a JSON echo of what it received. `POST /orders` records an order, `GET /orders` lists them and `DELETE /orders` forgets them. `PUT /control` makes it misbehave while it runs: `fail_health`, `fail_rate` (answer 500), `crash_rate` (drop the connection), `delay`. `relayctl fail` and `relayctl heal` are shortcuts for it.

## Tests

```
go test -race ./...
```
