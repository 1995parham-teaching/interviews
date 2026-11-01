package proxy

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"syscall"
)

// A request may be sent to a second backend only when that cannot do the
// work twice. That is true of an idempotent method whatever happened, and of
// any method when the first backend never received the request.

// idempotent methods can be repeated without changing the result.
func idempotent(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	}
	return false
}

// shouldRetry decides whether the outcome of one attempt is worth another.
func shouldRetry(r *http.Request, resp *http.Response, err error) bool {
	if r.Context().Err() != nil {
		return false // the client is gone
	}
	if err != nil {
		return idempotent(r.Method) || notSent(err)
	}
	switch resp.StatusCode {
	case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return idempotent(r.Method)
	}
	return false
}

// notSent reports whether err means the backend never received the request:
//
//   - the connection was refused, so nothing was written;
//   - the backend had already closed a kept-alive connection when the
//     request was written to it, which surfaces as EOF or a reset. The
//     backend closed the connection before reading anything from it.
func notSent(err error) bool {
	var op *net.OpError
	if errors.As(err, &op) && op.Op == "dial" {
		return true
	}
	return errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, syscall.ECONNRESET)
}

// timedOut reports whether err is the per-attempt timeout.
func timedOut(err error) bool {
	var ne net.Error
	return errors.Is(err, context.DeadlineExceeded) || (errors.As(err, &ne) && ne.Timeout())
}
