# Relay

Relay is a small distributed HTTP load balancer. Several Relay nodes run side by side, each one able to take any client request. They agree on the list of backends by gossip, and they route with a consistent-hash ring, so every node sends the same client to the same backend.

Only the Go standard library is used. Go 1.27 or later.

## How it works

- **Ring.** Every active backend is placed on a consistent-hash ring at 64 points. A request goes to the backend that owns its affinity key: the `X-Relay-Key` header if the client sends one, otherwise the client's IP address. If that backend is down, the request goes to the next backend on the ring.
- **Health.** Each node checks every backend on its own with `GET /health`, once a second. A backend is marked down after 3 failed checks in a row, and up again after 2 successful checks in a row. Health is not shared between nodes.
- **Retries.** If a backend cannot be reached or answers with a 5xx, the request is tried once more on the next backend.
- **Backend list.** All nodes start with the same seed list. An operator can add, drain or remove a backend through any node; the change carries a version and reaches the other nodes by gossip within a few seconds. For each backend the entry with the higher version wins.

```
            ┌──────────┐   gossip   ┌──────────┐   gossip   ┌──────────┐
 clients ──▶│ relay-1  │◀──────────▶│ relay-2  │◀──────────▶│ relay-3  │
            └────┬─────┘            └────┬─────┘            └────┬─────┘
                 │      ring + health    │                       │
                 ▼                       ▼                       ▼
            backend :9001          backend :9002           backend :9003
```

## Layout

```
cmd/relay      one Relay node
cmd/backend    a demo backend to put behind it
ring           the consistent-hash ring
health         active health checks
cluster        the shared backend list and its versions
proxy          forwarding, affinity and retries
node           wires the above together; admin API and gossip
```

## Running it

```
examples/cluster.sh
```

starts three backends on ports 9001 to 9003 and three nodes, `relay-1` to `relay-3`, on ports 8001 to 8003. Logs go to `logs/`. Then:

```
curl -s -H 'X-Relay-Key: alice' localhost:8001/hello
curl -s -H 'X-Relay-Key: alice' localhost:8002/hello      # same backend
curl -s localhost:8001/relay/status
```

A single node by hand:

```
go run ./cmd/relay -id relay-1 -listen 127.0.0.1:8001 \
    -peers 127.0.0.1:8002,127.0.0.1:8003 \
    -backends 127.0.0.1:9001,127.0.0.1:9002,127.0.0.1:9003
```

## Node API

Everything under `/relay/` is the node's own API; every other path is forwarded to a backend.

| Request | Effect |
| --- | --- |
| `GET /relay/status` | The node's id, its peers, the backends on its ring, and every known backend with its status, version and health as this node sees it |
| `PUT /relay/backends/{host:port}` with `{"status":"active"}` or `{"status":"draining"}` | Adds the backend or changes its status |
| `DELETE /relay/backends/{host:port}` | Removes the backend |
| `POST /relay/gossip` | Used between nodes: takes a list of entries, merges it, answers with the node's own list |

Responses that went through Relay carry `X-Relay-Node` and `X-Relay-Backend`.

## Demo backend

`cmd/backend` answers every request with a JSON echo of what it received, including `X-Forwarded-For` and `X-Relay-Key`. `POST /orders` records an order and `GET /orders` lists them. It can play a failing service while it runs:

```
curl -X PUT localhost:9002/control -d '{"fail_health": true}'     # /health answers 503
curl -X PUT localhost:9002/control -d '{"fail_rate": 0.5}'        # half of the requests answer 500
curl -X PUT localhost:9002/control -d '{"delay": "300ms"}'
curl -X PUT localhost:9002/control -d '{}'                        # back to normal
```

## Tests

```
go test -race ./...
```
