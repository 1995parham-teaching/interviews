// Command relay runs one load balancer node.
package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"relay/node"
)

func main() {
	id := flag.String("id", "", "node id, shown in X-Relay-Node (default: the listen address)")
	listen := flag.String("listen", "127.0.0.1:8001", "address to serve on")
	peers := flag.String("peers", "", "other nodes, comma separated host:port")
	backends := flag.String("backends", "", "seed backends, comma separated host:port; the same on every node")
	check := flag.Duration("check-interval", time.Second, "time between health checks")
	gossip := flag.Duration("gossip-interval", time.Second, "time between gossip rounds")
	flag.Parse()

	if *id == "" {
		*id = *listen
	}
	n := node.New(node.Config{
		ID:             *id,
		Peers:          split(*peers),
		Backends:       split(*backends),
		CheckInterval:  *check,
		GossipInterval: *gossip,
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{Addr: *listen, Handler: n.Handler(), ReadHeaderTimeout: 5 * time.Second}
	go n.Run(ctx)
	go func() {
		<-ctx.Done()
		shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(shutdown)
	}()

	slog.Info("relay up", "node", *id, "listen", *listen, "peers", *peers, "backends", *backends)
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		slog.Error("relay stopped", "err", err)
		os.Exit(1)
	}
}

func split(s string) []string {
	var out []string
	for part := range strings.SplitSeq(s, ",") {
		if part = strings.TrimSpace(part); part != "" {
			out = append(out, part)
		}
	}
	return out
}
