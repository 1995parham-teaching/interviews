// Command relayctl pokes a running Relay cluster: it shows every node's view
// side by side, changes backends through a chosen node, makes a demo backend
// misbehave, and sends batches of requests.
//
// It finds the cluster in logs/cluster.json, written by `go run ./cmd/cluster`;
// without that file it assumes nodes on 8001-8003 and backends on 9001-9003.
// Nodes can be named relay-1, 1 or 8001; backends backend-9001 or 9001.
package main

import (
	"bytes"
	"cmp"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"maps"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

const usage = `relayctl: poke a running Relay cluster

  status                          every backend as each node sees it, and each node's runtime
  activate <node> <backend>       PUT {"status":"active"} through <node>; adds an unknown backend
  drain    <node> <backend>       PUT {"status":"draining"} through <node>
  weight   <node> <backend> <w>   PUT {"weight":w} through <node>
  remove   <node> <backend>       DELETE through <node>

  fail <backend> health           its /health answers 503
  fail <backend> rate <0..1>      that share of its requests answer 500 (after doing the work)
  fail <backend> crash <0..1>     that share of its requests drop the connection (after doing the work)
  fail <backend> delay <dur>      every request waits, e.g. 300ms
  heal <backend>                  back to normal

  send [flags] [path]             send requests and summarise where they went
      -node relay-1 | all   -n 20   -method GET   -body ''   -v
      -key 'user-%d'  the affinity key, %d is the request number; -key '' sends none
      %d in the path is replaced too; every request opens a new connection

  orders                          the orders each backend has recorded, by key
  orders reset                    forget all orders on every backend

Nodes: relay-1, 1 or 8001. Backends: backend-9001 or 9001.
`

type layout struct {
	Nodes    map[string]string `json:"nodes"`
	Backends map[string]string `json:"backends"`
}

var (
	cl     layout
	client = &http.Client{Timeout: 5 * time.Second}
)

func main() {
	cl = load()
	args := os.Args[1:]
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		io.WriteString(os.Stdout, usage)
		return
	}
	var err error
	switch cmd, rest := args[0], args[1:]; cmd {
	case "status":
		err = status()
	case "activate", "drain":
		err = need(rest, 2, func() error {
			st := map[string]string{"activate": "active", "drain": "draining"}[cmd]
			return put(rest[0], rest[1], fmt.Sprintf(`{"status":%q}`, st))
		})
	case "weight":
		err = need(rest, 3, func() error { return put(rest[0], rest[1], `{"weight":`+rest[2]+`}`) })
	case "remove":
		err = need(rest, 2, func() error { return remove(rest[0], rest[1]) })
	case "fail":
		err = failCmd(rest)
	case "heal":
		err = need(rest, 1, func() error { return control(rest[0], `{}`) })
	case "send":
		err = send(rest)
	case "orders":
		if len(rest) == 1 && rest[0] == "reset" {
			err = resetOrders()
		} else {
			err = orders()
		}
	default:
		err = fmt.Errorf("unknown command %q; run relayctl help", cmd)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "relayctl:", err)
		os.Exit(1)
	}
}

func need(args []string, n int, f func() error) error {
	if len(args) != n {
		return fmt.Errorf("wrong number of arguments; run relayctl help")
	}
	return f()
}

// ---------------------------------------------------------------- names

func load() layout {
	l := layout{Nodes: map[string]string{}, Backends: map[string]string{}}
	if b, err := os.ReadFile(filepath.Join("logs", "cluster.json")); err == nil && json.Unmarshal(b, &l) == nil {
		return l
	}
	for i := range 3 {
		l.Nodes[fmt.Sprintf("relay-%d", i+1)] = fmt.Sprintf("127.0.0.1:%d", 8001+i)
		l.Backends[fmt.Sprintf("backend-%d", 9001+i)] = fmt.Sprintf("127.0.0.1:%d", 9001+i)
	}
	return l
}

func resolve(m map[string]string, kind, name string) (string, error) {
	if a, ok := m[name]; ok {
		return a, nil
	}
	for n, a := range m {
		if strings.HasSuffix(a, ":"+name) || a == name || strings.TrimPrefix(n, "relay-") == name {
			return a, nil
		}
	}
	if strings.Contains(name, ":") { // any ip:port, e.g. a backend that is not running yet
		return name, nil
	}
	return "", fmt.Errorf("unknown %s %q (known: %s)", kind, name, strings.Join(slices.Sorted(maps.Keys(m)), ", "))
}

func node(name string) (string, error)    { return resolve(cl.Nodes, "node", name) }
func backend(name string) (string, error) { return resolve(cl.Backends, "backend", name) }

// short turns 127.0.0.1:9002 into :9002 for tables.
func short(addr string) string {
	if i := strings.LastIndex(addr, ":"); i >= 0 {
		return addr[i:]
	}
	return addr
}

// ---------------------------------------------------------------- status

type entry struct {
	Backend string `json:"backend"`
	Status  string `json:"status"`
	Weight  int    `json:"weight"`
	Version uint64 `json:"version"`
	Healthy bool   `json:"healthy"`
}

type nodeStatus struct {
	Node  string `json:"node"`
	Clock uint64 `json:"clock"`
	Ring  []struct {
		Backend string `json:"backend"`
	} `json:"ring"`
	Backends []entry `json:"backends"`
	Pins     int     `json:"pins"`
	Runtime  struct {
		Goroutines int     `json:"goroutines"`
		HeapMB     float64 `json:"heap_mb"`
	} `json:"runtime"`
}

func status() error {
	names := slices.Sorted(maps.Keys(cl.Nodes))
	views := map[string]*nodeStatus{}
	all := map[string]bool{}
	for _, n := range names {
		var s nodeStatus
		if err := getJSON("http://"+cl.Nodes[n]+"/relay/status", &s); err != nil {
			continue
		}
		views[n] = &s
		for _, e := range s.Backends {
			all[e.Backend] = true
		}
	}
	if len(views) == 0 {
		return fmt.Errorf("no node answers at %s; start the cluster with: go run ./cmd/cluster", strings.Join(slices.Sorted(maps.Values(cl.Nodes)), ", "))
	}

	const w = 24
	row := func(label string, cell func(v *nodeStatus) string) {
		fmt.Printf("%-11s", label)
		for _, n := range names {
			c := "unreachable"
			if v := views[n]; v != nil {
				c = cell(v)
			}
			fmt.Printf("%-*s", w, c)
		}
		fmt.Println()
	}
	row("", func(v *nodeStatus) string { return v.Node })
	for _, b := range slices.Sorted(maps.Keys(all)) {
		row(short(b), func(v *nodeStatus) string {
			for _, e := range v.Backends {
				if e.Backend == b {
					if e.Status == "removed" {
						return fmt.Sprintf("removed v%d", e.Version)
					}
					h := "up"
					if !e.Healthy {
						h = "DOWN"
					}
					return fmt.Sprintf("%s w%d v%d %s", e.Status, e.Weight, e.Version, h)
				}
			}
			return "—"
		})
	}
	row("ring", func(v *nodeStatus) string {
		var r []string
		for _, m := range v.Ring {
			r = append(r, short(m.Backend))
		}
		if len(r) == 0 {
			return "(empty)"
		}
		return strings.Join(r, " ")
	})
	row("clock", func(v *nodeStatus) string { return fmt.Sprint(v.Clock) })
	row("pins", func(v *nodeStatus) string { return fmt.Sprint(v.Pins) })
	row("goroutines", func(v *nodeStatus) string { return fmt.Sprint(v.Runtime.Goroutines) })
	row("heap MB", func(v *nodeStatus) string { return fmt.Sprint(v.Runtime.HeapMB) })
	fmt.Println("\nstatus wWEIGHT vVERSION health, as each node sees it; — means the node does not know the backend.")
	return nil
}

// ---------------------------------------------------------------- changes

func put(nodeName, backendName, body string) error {
	n, err := node(nodeName)
	if err != nil {
		return err
	}
	b, err := backend(backendName)
	if err != nil {
		return err
	}
	out, code, err := do("PUT", "http://"+n+"/relay/backends/"+b, body, nil)
	if err != nil {
		return err
	}
	fmt.Printf("%s %d %s\n", nodeName, code, strings.TrimSpace(compact(out)))
	return nil
}

func remove(nodeName, backendName string) error {
	n, err := node(nodeName)
	if err != nil {
		return err
	}
	b, err := backend(backendName)
	if err != nil {
		return err
	}
	out, code, err := do("DELETE", "http://"+n+"/relay/backends/"+b, "", nil)
	if err != nil {
		return err
	}
	fmt.Printf("%s %d %s\n", nodeName, code, strings.TrimSpace(compact(out)))
	return nil
}

func failCmd(args []string) error {
	if len(args) >= 2 && args[1] == "health" {
		return control(args[0], `{"fail_health": true}`)
	}
	if len(args) != 3 {
		return fmt.Errorf("usage: fail <backend> health | rate <0..1> | crash <0..1> | delay <dur>")
	}
	switch args[1] {
	case "rate":
		return control(args[0], `{"fail_rate": `+args[2]+`}`)
	case "crash":
		return control(args[0], `{"crash_rate": `+args[2]+`}`)
	case "delay":
		return control(args[0], fmt.Sprintf(`{"delay": %q}`, args[2]))
	}
	return fmt.Errorf("unknown failure %q: health, rate, crash or delay", args[1])
}

func control(backendName, body string) error {
	b, err := backend(backendName)
	if err != nil {
		return err
	}
	out, code, err := do("PUT", "http://"+b+"/control", body, nil)
	if err != nil {
		return err
	}
	if code != http.StatusOK {
		return fmt.Errorf("%s answered %d: %s", short(b), code, strings.TrimSpace(out))
	}
	fmt.Printf("%s now %s\n", short(b), strings.TrimSpace(compact(out)))
	return nil
}

// ---------------------------------------------------------------- send

func send(args []string) error {
	fs := flag.NewFlagSet("send", flag.ContinueOnError)
	nodeName := fs.String("node", "relay-1", "node to send to, or all")
	n := fs.Int("n", 20, "number of requests")
	method := fs.String("method", "GET", "HTTP method")
	body := fs.String("body", "", "request body")
	key := fs.String("key", "user-%d", "affinity key; %d is the request number; empty sends none")
	verbose := fs.Bool("v", false, "print every request")
	if err := fs.Parse(args); err != nil {
		return err
	}
	path := "/hello"
	if fs.NArg() > 0 {
		path = fs.Arg(0)
	}
	names := []string{*nodeName}
	if *nodeName == "all" {
		names = slices.Sorted(maps.Keys(cl.Nodes))
	}
	for _, name := range names {
		addr, err := node(name)
		if err != nil {
			return err
		}
		if _, _, err := do("GET", "http://"+addr+"/relay/status", "", nil); err != nil {
			return fmt.Errorf("%s does not answer at %s; is the cluster running? (go run ./cmd/cluster)", name, addr)
		}
		sendTo(name, addr, *n, *method, path, *body, *key, *verbose)
	}
	return nil
}

func fill(pattern string, i int) string {
	if strings.Contains(pattern, "%d") {
		return fmt.Sprintf(pattern, i)
	}
	return pattern
}

func sendTo(name, addr string, n int, method, path, body, key string, verbose bool) {
	// A new connection per request, as separate clients (or curl in a loop) would.
	c := &http.Client{Timeout: 10 * time.Second, Transport: &http.Transport{DisableKeepAlives: true}}
	byBackend := map[string]int{}
	byCode := map[int]int{}
	var failed []string
	for i := range n {
		k := ""
		if key != "" {
			k = fill(key, i)
		}
		h := http.Header{}
		if k != "" {
			h.Set("X-Relay-Key", k)
		}
		_, code, be, err := doWith(c, method, "http://"+addr+fill(path, i), body, h)
		label := k
		if label == "" {
			label = fmt.Sprintf("#%d", i)
		}
		if err != nil {
			byCode[0]++
			failed = append(failed, label+" error")
			continue
		}
		byCode[code]++
		if be == "" {
			be = "(none)"
		}
		byBackend[short(be)]++
		if code >= 300 {
			failed = append(failed, fmt.Sprintf("%s %d", label, code))
		}
		if verbose {
			fmt.Printf("  %-12s -> %-8s %d\n", label, short(be), code)
		}
	}

	fmt.Printf("%d x %s %s via %s\n", n, method, path, name)
	fmt.Print("  status: ")
	for _, code := range slices.Sorted(maps.Keys(byCode)) {
		label := fmt.Sprint(code)
		if code == 0 {
			label = "error"
		}
		fmt.Printf("%s x%d  ", label, byCode[code])
	}
	fmt.Print("\n  served by: ")
	for _, b := range slices.Sorted(maps.Keys(byBackend)) {
		fmt.Printf("%s x%d  ", b, byBackend[b])
	}
	fmt.Println()
	if len(failed) > 0 {
		show := failed[:min(len(failed), 12)]
		fmt.Printf("  not 2xx: %s", strings.Join(show, ", "))
		if len(failed) > len(show) {
			fmt.Printf(" ... and %d more", len(failed)-len(show))
		}
		fmt.Println()
	}
}

// ---------------------------------------------------------------- orders

type order struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
	Key  string `json:"key"`
}

func orders() error {
	type row struct {
		backend string
		o       order
	}
	var rows []row
	for _, name := range slices.Sorted(maps.Keys(cl.Backends)) {
		addr := cl.Backends[name]
		var r struct {
			Orders []order `json:"orders"`
		}
		if err := getJSON("http://"+addr+"/orders", &r); err != nil {
			fmt.Printf("%s unreachable\n", short(addr))
			continue
		}
		fmt.Printf("%s %d orders\n", short(addr), len(r.Orders))
		for _, o := range r.Orders {
			rows = append(rows, row{short(addr), o})
		}
	}
	slices.SortStableFunc(rows, func(a, b row) int { return cmp.Compare(a.o.Key, b.o.Key) })
	if len(rows) > 0 {
		fmt.Printf("\n  %-12s %-8s %s\n", "key", "backend", "body")
		for _, r := range rows {
			fmt.Printf("  %-12s %-8s %s\n", r.o.Key, r.backend, r.o.Body)
		}
	}
	return nil
}

func resetOrders() error {
	for _, name := range slices.Sorted(maps.Keys(cl.Backends)) {
		addr := cl.Backends[name]
		if _, code, err := do("DELETE", "http://"+addr+"/orders", "", nil); err != nil || code != http.StatusNoContent {
			fmt.Printf("%s not reset (%v %d)\n", short(addr), err, code)
			continue
		}
		fmt.Printf("%s 0 orders\n", short(addr))
	}
	return nil
}

// ---------------------------------------------------------------- http

func getJSON(url string, v any) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(v)
}

func do(method, url, body string, h http.Header) (string, int, error) {
	out, code, _, err := doWith(client, method, url, body, h)
	return out, code, err
}

func doWith(c *http.Client, method, url, body string, h http.Header) (string, int, string, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return "", 0, "", err
	}
	maps.Copy(req.Header, h)
	resp, err := c.Do(req)
	if err != nil {
		return "", 0, "", err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return string(b), resp.StatusCode, resp.Header.Get("X-Relay-Backend"), nil
}

func compact(s string) string {
	var buf bytes.Buffer
	if json.Compact(&buf, []byte(s)) == nil {
		return buf.String()
	}
	return s
}
