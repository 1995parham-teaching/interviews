// Command relay runs one load balancer node.
//
//	relay -id relay-1 -listen 127.0.0.1:8001 \
//	    -peers 127.0.0.1:8002,127.0.0.1:8003 \
//	    -backends 127.0.0.1:9001,127.0.0.1:9002=2
//
// Settings can also come from a JSON file with -config; flags given on the
// command line override it. Run with -h for every flag.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"relay/config"
	"relay/node"
)

func main() {
	cfg, err := config.Load(os.Args[1:])
	if errors.Is(err, flag.ErrHelp) {
		return
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "relay:", err)
		os.Exit(2)
	}
	n, err := node.New(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "relay:", err)
		os.Exit(2)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{Addr: cfg.Listen, Handler: n.Handler(), ReadHeaderTimeout: 5 * time.Second}
	go n.Run(ctx)
	go func() {
		<-ctx.Done()
		shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(shutdown)
	}()

	slog.Info("relay up", "node", cfg.ID, "listen", cfg.Listen, "peers", cfg.Peers, "backends", len(cfg.Backends))
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		slog.Error("relay stopped", "err", err)
		os.Exit(1)
	}
}
