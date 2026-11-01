// Command cluster runs a whole Relay cluster on this machine: three backends
// and three nodes, each in its own process, with logs in ./logs.
//
//	go run ./cmd/cluster
//
// While it runs, it reads commands from the terminal:
//
//	ps                    list the processes and whether they are running
//	stop  <name>          stop one process (relay-2, backend-9003, ...)
//	start <name>          start it again, with the same flags
//	restart <name>        both
//	logs <name>           print the last lines of its log
//	quit                  stop everything (Ctrl-C does the same)
//
// The addresses are written to logs/cluster.json, where relayctl finds them.
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
	"syscall"
	"time"
)

type proc struct {
	name string
	bin  string
	args []string

	mu   sync.Mutex
	cmd  *exec.Cmd
	done chan struct{} // closed when cmd has exited
}

// Layout is written to logs/cluster.json for relayctl.
type Layout struct {
	Nodes    map[string]string `json:"nodes"`    // relay-1 -> 127.0.0.1:8001
	Backends map[string]string `json:"backends"` // backend-9001 -> 127.0.0.1:9001
}

func main() {
	nodes := flag.Int("nodes", 3, "number of Relay nodes")
	backends := flag.Int("backends", 3, "number of backends")
	nodePort := flag.Int("node-port", 8001, "port of the first node; the others follow")
	backendPort := flag.Int("backend-port", 9001, "port of the first backend; the others follow")
	tombstoneTTL := flag.Duration("tombstone-ttl", 30*time.Second, "passed to every node; short, so expiry can be watched")
	pinTTL := flag.Duration("pin-ttl", 10*time.Minute, "passed to every node")
	trusted := flag.String("trusted-proxies", "127.0.0.1/32", "passed to every node: CIDRs whose X-Forwarded-For is believed")
	flag.Parse()

	exe := func(name string) string {
		if runtime.GOOS == "windows" {
			name += ".exe"
		}
		return filepath.Join("bin", name)
	}
	fmt.Println("building bin/relay, bin/backend and bin/relayctl ...")
	build := exec.Command("go", "build", "-o", "bin"+string(filepath.Separator), "./cmd/relay", "./cmd/backend", "./cmd/relayctl")
	build.Stdout, build.Stderr = os.Stdout, os.Stderr
	if err := build.Run(); err != nil {
		fail("build failed: %v", err)
	}
	if err := os.MkdirAll("logs", 0o755); err != nil {
		fail("%v", err)
	}

	layout := Layout{Nodes: map[string]string{}, Backends: map[string]string{}}
	var procs []*proc
	var backendAddrs, nodeAddrs []string
	for i := range *backends {
		addr := fmt.Sprintf("127.0.0.1:%d", *backendPort+i)
		name := fmt.Sprintf("backend-%d", *backendPort+i)
		backendAddrs = append(backendAddrs, addr)
		layout.Backends[name] = addr
		procs = append(procs, &proc{name: name, bin: exe("backend"), args: []string{"-listen", addr}})
	}
	for i := range *nodes {
		nodeAddrs = append(nodeAddrs, fmt.Sprintf("127.0.0.1:%d", *nodePort+i))
	}
	for i, addr := range nodeAddrs {
		name := fmt.Sprintf("relay-%d", i+1)
		var peers []string
		for j, p := range nodeAddrs {
			if j != i {
				peers = append(peers, p)
			}
		}
		layout.Nodes[name] = addr
		procs = append(procs, &proc{name: name, bin: exe("relay"), args: []string{
			"-id", name, "-listen", addr,
			"-peers", strings.Join(peers, ","),
			"-backends", strings.Join(backendAddrs, ","),
			"-tombstone-ttl", tombstoneTTL.String(),
			"-pin-ttl", pinTTL.String(),
			"-trusted-proxies", *trusted,
		}})
	}

	for _, addr := range slices.Concat(backendAddrs, nodeAddrs) {
		if l, err := net.Listen("tcp", addr); err != nil {
			fail("%s is in use; stop whatever runs there, or pick other ports with -node-port and -backend-port", addr)
		} else {
			l.Close()
		}
	}

	b, _ := json.MarshalIndent(layout, "", "  ")
	if err := os.WriteFile(filepath.Join("logs", "cluster.json"), b, 0o644); err != nil {
		fail("%v", err)
	}

	for _, p := range procs {
		if err := p.start(); err != nil {
			stopAll(procs)
			fail("starting %s: %v", p.name, err)
		}
	}
	waitReady(backendAddrs, nodeAddrs)

	fmt.Println()
	for _, p := range procs {
		fmt.Printf("  %-14s %s\n", p.name, addrOf(p))
	}
	fmt.Println()
	fmt.Println("cluster up. Logs in ./logs. In another terminal:")
	fmt.Println("  bin/relayctl status          (or: go run ./cmd/relayctl status)")
	fmt.Println("  curl -H 'X-Relay-Key: alice' localhost:8001/hello")
	fmt.Println("Type 'help' for commands here; Ctrl-C or 'quit' stops everything.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	lines := make(chan string)
	go func() {
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			lines <- sc.Text()
		}
		close(lines)
	}()

	byName := map[string]*proc{}
	for _, p := range procs {
		byName[p.name] = p
	}
	fmt.Print("> ")
	for {
		select {
		case <-sig:
			fmt.Println("\nstopping ...")
			stopAll(procs)
			return
		case line, ok := <-lines:
			if !ok { // stdin closed: keep running until a signal
				lines = nil
				continue
			}
			if quit := command(strings.Fields(line), procs, byName); quit {
				stopAll(procs)
				return
			}
			fmt.Print("> ")
		}
	}
}

func command(args []string, procs []*proc, byName map[string]*proc) (quit bool) {
	if len(args) == 0 {
		return false
	}
	need := func() *proc {
		if len(args) != 2 {
			fmt.Println("usage:", args[0], "<name>   (names: see 'ps')")
			return nil
		}
		name := args[1]
		if p, ok := byName[name]; ok {
			return p
		}
		// accept 8002, 9003 or relay2 as shorthands
		for _, p := range procs {
			if strings.HasSuffix(addrOf(p), ":"+name) || strings.ReplaceAll(p.name, "-", "") == name {
				return p
			}
		}
		fmt.Println("unknown name", name, "(see 'ps')")
		return nil
	}
	switch args[0] {
	case "help", "?":
		fmt.Println("ps | stop <name> | start <name> | restart <name> | logs <name> | quit")
	case "ps":
		for _, p := range procs {
			state := "stopped"
			if p.running() {
				state = "running"
			}
			fmt.Printf("  %-14s %-16s %s\n", p.name, addrOf(p), state)
		}
	case "stop":
		if p := need(); p != nil {
			p.stop()
			fmt.Println("stopped", p.name)
		}
	case "start":
		if p := need(); p != nil {
			if p.running() {
				fmt.Println(p.name, "is already running")
			} else if err := p.start(); err != nil {
				fmt.Println("error:", err)
			} else {
				fmt.Println("started", p.name)
			}
		}
	case "restart":
		if p := need(); p != nil {
			p.stop()
			if err := p.start(); err != nil {
				fmt.Println("error:", err)
			} else {
				fmt.Println("restarted", p.name)
			}
		}
	case "logs":
		if p := need(); p != nil {
			tail(filepath.Join("logs", p.name+".log"), 15)
		}
	case "quit", "exit":
		return true
	default:
		fmt.Println("unknown command; type 'help'")
	}
	return false
}

func (p *proc) start() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	f, err := os.OpenFile(filepath.Join("logs", p.name+".log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	cmd := exec.Command(p.bin, p.args...)
	cmd.Stdout, cmd.Stderr = f, f
	if err := cmd.Start(); err != nil {
		f.Close()
		return err
	}
	done := make(chan struct{})
	p.cmd, p.done = cmd, done
	go func() {
		cmd.Wait()
		f.Close()
		close(done)
	}()
	return nil
}

func (p *proc) stop() {
	p.mu.Lock()
	cmd, done := p.cmd, p.done
	p.cmd, p.done = nil, nil
	p.mu.Unlock()
	if cmd == nil {
		return
	}
	if runtime.GOOS == "windows" {
		cmd.Process.Kill()
	} else {
		cmd.Process.Signal(syscall.SIGTERM)
	}
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		cmd.Process.Kill()
	}
}

func (p *proc) running() bool {
	p.mu.Lock()
	done := p.done
	p.mu.Unlock()
	if done == nil {
		return false
	}
	select {
	case <-done:
		return false
	default:
		return true
	}
}

func stopAll(procs []*proc) {
	var wg sync.WaitGroup
	for _, p := range procs {
		wg.Go(p.stop)
	}
	wg.Wait()
}

func addrOf(p *proc) string {
	for i, a := range p.args {
		if a == "-listen" && i+1 < len(p.args) {
			return p.args[i+1]
		}
	}
	return ""
}

func waitReady(backends, nodes []string) {
	deadline := time.Now().Add(10 * time.Second)
	for _, addr := range slices.Concat(backends, nodes) {
		for {
			c, err := net.DialTimeout("tcp", addr, 200*time.Millisecond)
			if err == nil {
				c.Close()
				break
			}
			if time.Now().After(deadline) {
				fmt.Printf("warning: %s is not answering yet; see logs/\n", addr)
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func tail(path string, n int) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer f.Close()
	b, _ := io.ReadAll(f)
	lines := strings.Split(strings.TrimRight(string(b), "\n"), "\n")
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	for _, l := range lines {
		fmt.Println("  " + l)
	}
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "cluster: "+format+"\n", args...)
	os.Exit(1)
}
