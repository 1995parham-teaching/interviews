// Package config holds a node's settings: defaults, an optional JSON file,
// then command-line flags, each overriding the one before.
package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/netip"
	"os"
	"strconv"
	"strings"
	"time"
)

// Backend is a seed backend with its weight.
type Backend struct {
	Addr   string `json:"addr"`
	Weight int    `json:"weight"`
}

// Config is everything a node needs to start.
type Config struct {
	ID     string `json:"id"`
	Listen string `json:"listen"`

	// Peers are the other nodes; every node gossips with every peer.
	Peers []string `json:"peers"`
	// Backends is the seed list. Every node is started with the same one.
	Backends []Backend `json:"backends"`

	CheckInterval Duration `json:"check_interval"`
	CheckTimeout  Duration `json:"check_timeout"`
	Fall          int      `json:"fall"`
	Rise          int      `json:"rise"`

	GossipInterval Duration `json:"gossip_interval"`
	GossipTimeout  Duration `json:"gossip_timeout"`
	// TombstoneTTL is how long a removed backend is remembered as removed.
	TombstoneTTL Duration `json:"tombstone_ttl"`

	Retries        int      `json:"retries"`
	TryTimeout     Duration `json:"try_timeout"`
	PinTTL         Duration `json:"pin_ttl"`
	MaxBody        int64    `json:"max_body"`
	TrustedProxies []string `json:"trusted_proxies"`

	// AdminToken, if set, is required as a bearer token on admin writes.
	AdminToken string `json:"admin_token"`
	// ClusterSecret, if set, must be sent by peers on gossip.
	ClusterSecret string `json:"cluster_secret"`
}

// Default returns the settings used when nothing else is given.
func Default() Config {
	return Config{
		Listen:         "127.0.0.1:8001",
		CheckInterval:  Duration(time.Second),
		CheckTimeout:   Duration(time.Second),
		Fall:           3,
		Rise:           2,
		GossipInterval: Duration(time.Second),
		GossipTimeout:  Duration(2 * time.Second),
		TombstoneTTL:   Duration(10 * time.Minute),
		Retries:        1,
		TryTimeout:     Duration(5 * time.Second),
		PinTTL:         Duration(10 * time.Minute),
		MaxBody:        1 << 20,
	}
}

// Load reads the configuration: defaults, then the file named by -config if
// any, then the other flags that were set explicitly.
func Load(args []string) (Config, error) {
	c := Default()
	fs := flag.NewFlagSet("relay", flag.ContinueOnError)
	file := fs.String("config", "", "JSON configuration file")
	id := fs.String("id", "", "node id (default: the listen address)")
	listen := fs.String("listen", c.Listen, "address to serve on")
	peers := fs.String("peers", "", "other nodes, comma separated host:port")
	backends := fs.String("backends", "", "seed backends, comma separated host:port[=weight]")
	check := fs.Duration("check-interval", time.Duration(c.CheckInterval), "time between health checks")
	gossip := fs.Duration("gossip-interval", time.Duration(c.GossipInterval), "time between gossip rounds")
	ttl := fs.Duration("tombstone-ttl", time.Duration(c.TombstoneTTL), "how long a removed backend is remembered")
	pin := fs.Duration("pin-ttl", time.Duration(c.PinTTL), "how long an idle session pin is kept")
	retries := fs.Int("retries", c.Retries, "extra attempts after the first")
	trusted := fs.String("trusted-proxies", "", "comma separated CIDRs whose X-Forwarded-For is believed")
	token := fs.String("admin-token", "", "bearer token required on admin writes")
	secret := fs.String("cluster-secret", "", "secret every node sends and requires on gossip")
	if err := fs.Parse(args); err != nil {
		return c, err
	}

	if *file != "" {
		b, err := os.ReadFile(*file)
		if err != nil {
			return c, err
		}
		if err := json.Unmarshal(b, &c); err != nil {
			return c, fmt.Errorf("%s: %w", *file, err)
		}
	}

	set := map[string]bool{}
	fs.Visit(func(f *flag.Flag) { set[f.Name] = true })
	if set["id"] {
		c.ID = *id
	}
	if set["listen"] {
		c.Listen = *listen
	}
	if set["peers"] {
		c.Peers = splitList(*peers)
	}
	if set["backends"] {
		bs, err := ParseBackends(*backends)
		if err != nil {
			return c, err
		}
		c.Backends = bs
	}
	if set["check-interval"] {
		c.CheckInterval = Duration(*check)
	}
	if set["gossip-interval"] {
		c.GossipInterval = Duration(*gossip)
	}
	if set["tombstone-ttl"] {
		c.TombstoneTTL = Duration(*ttl)
	}
	if set["pin-ttl"] {
		c.PinTTL = Duration(*pin)
	}
	if set["retries"] {
		c.Retries = *retries
	}
	if set["trusted-proxies"] {
		c.TrustedProxies = splitList(*trusted)
	}
	if set["admin-token"] {
		c.AdminToken = *token
	}
	if set["cluster-secret"] {
		c.ClusterSecret = *secret
	}
	if c.ID == "" {
		c.ID = c.Listen
	}
	return c, c.Validate()
}

// Validate reports the first setting that cannot work.
func (c Config) Validate() error {
	var errs []error
	if c.CheckInterval <= 0 || c.GossipInterval <= 0 {
		errs = append(errs, errors.New("intervals must be positive"))
	}
	if c.CheckTimeout <= 0 || c.GossipTimeout <= 0 || c.TryTimeout <= 0 {
		errs = append(errs, errors.New("timeouts must be positive"))
	}
	if c.Fall < 1 || c.Rise < 1 {
		errs = append(errs, errors.New("fall and rise must be at least 1"))
	}
	if c.TombstoneTTL <= 0 || c.PinTTL <= 0 {
		errs = append(errs, errors.New("tombstone_ttl and pin_ttl must be positive"))
	}
	if c.Retries < 0 {
		errs = append(errs, errors.New("retries must not be negative"))
	}
	if c.MaxBody <= 0 {
		errs = append(errs, errors.New("max_body must be positive"))
	}
	for _, b := range c.Backends {
		if _, err := netip.ParseAddrPort(b.Addr); err != nil {
			errs = append(errs, fmt.Errorf("backend %q: %w", b.Addr, err))
		}
		if b.Weight < 1 || b.Weight > MaxWeight {
			errs = append(errs, fmt.Errorf("backend %q: weight %d is not between 1 and %d", b.Addr, b.Weight, MaxWeight))
		}
	}
	if _, err := c.TrustedPrefixes(); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

// MaxWeight is the largest weight a backend may have.
const MaxWeight = 10

// TrustedPrefixes parses TrustedProxies.
func (c Config) TrustedPrefixes() ([]netip.Prefix, error) {
	var out []netip.Prefix
	for _, s := range c.TrustedProxies {
		p, err := netip.ParsePrefix(s)
		if err != nil {
			return nil, fmt.Errorf("trusted proxy %q: %w", s, err)
		}
		out = append(out, p.Masked())
	}
	return out, nil
}

// ParseBackends reads "host:port[=weight],..." with a default weight of 1.
func ParseBackends(s string) ([]Backend, error) {
	var out []Backend
	for _, part := range splitList(s) {
		addr, w, found := strings.Cut(part, "=")
		b := Backend{Addr: addr, Weight: 1}
		if found {
			n, err := strconv.Atoi(w)
			if err != nil {
				return nil, fmt.Errorf("backend %q: bad weight", part)
			}
			b.Weight = n
		}
		out = append(out, b)
	}
	return out, nil
}

func splitList(s string) []string {
	var out []string
	for part := range strings.SplitSeq(s, ",") {
		if part = strings.TrimSpace(part); part != "" {
			out = append(out, part)
		}
	}
	return out
}

// Duration reads "1.5s" from JSON and writes it back the same way.
type Duration time.Duration

func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("duration must be a string such as \"1s\": %w", err)
	}
	v, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*d = Duration(v)
	return nil
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

// D returns the value as a time.Duration.
func (d Duration) D() time.Duration { return time.Duration(d) }
