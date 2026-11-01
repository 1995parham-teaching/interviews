package proxy

import (
	"net"
	"net/http"
	"net/netip"
	"strings"
)

// ClientIP returns the address of the client that sent r.
//
// Normally that is the peer of the connection. When the peer is a trusted
// proxy (a corporate proxy or an L4 balancer in front of Relay), the peer is
// not the client: the proxy says who the client was in X-Forwarded-For, and
// the address it wrote there is used instead. Headers from anyone else are
// ignored, since any client can write them.
func ClientIP(r *http.Request, trusted []netip.Prefix) string {
	peer := hostOnly(r.RemoteAddr)
	if !isTrusted(peer, trusted) {
		return peer
	}
	for part := range strings.SplitSeq(r.Header.Get("X-Forwarded-For"), ",") {
		if ip, err := netip.ParseAddr(strings.TrimSpace(part)); err == nil {
			return ip.Unmap().String()
		}
	}
	return peer
}

// forwardedFor is the X-Forwarded-For to send upstream: the chain we
// received with the client's peer address appended.
func forwardedFor(r *http.Request) string {
	peer := hostOnly(r.RemoteAddr)
	if prior := r.Header.Get("X-Forwarded-For"); prior != "" {
		return prior + ", " + peer
	}
	return peer
}

func hostOnly(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		host = addr
	}
	if ip, err := netip.ParseAddr(host); err == nil {
		return ip.Unmap().String()
	}
	return host
}

func isTrusted(addr string, trusted []netip.Prefix) bool {
	ip, err := netip.ParseAddr(addr)
	if err != nil {
		return false
	}
	for _, p := range trusted {
		if p.Contains(ip) {
			return true
		}
	}
	return false
}
