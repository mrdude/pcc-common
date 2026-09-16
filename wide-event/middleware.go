package wideevt

import (
	"bufio"
	"errors"
	"math/rand/v2"
	"net"
	"net/http"
	"time"

	"github.com/mrdude/pcc-common/server"
)

// CreateMiddleware returns a mux.MiddlewareFunc that
// attaches an Event to the context, and writes it to the sink.
func CreateMiddleware(sink Sink, serverName string) server.OptionsFn {
	return server.WithGlobalMiddleware(
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				ctx := req.Context()
				evt := NewEvent()
				ctx = WithEvent(ctx, evt)
				req = req.WithContext(ctx)

				start := time.Now()
				next.ServeHTTP(w, req)
				evt.Duration = time.Since(start)

				if shouldSampleEvent(req, evt) {
					sink.Write(evt)
				}
			})
		},
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				ctx := req.Context()
				evt := GetEvent(ctx)

				evt.ServerName = serverName

				evt.RequestMethod = req.Method
				evt.RequestPath = req.URL.Path
				evt.RequestURL = req.URL.String()

				wrap := &responseWrapper{delegate: w}
				next.ServeHTTP(wrap, req)

				// response fields
				evt.StatusCode = wrap.statusCode
			})
		},
	)
}

func shouldSampleEvent(req *http.Request, evt *Event) bool {
	// include all failed requests
	if evt.StatusCode >= 400 || evt.StatusCode <= 599 {
		return true
	}

	// include all throttled events
	if _, ok := evt.data["rate_limited"]; ok {
		return true
	}

	// include slow requests
	if evt.Duration >= 9*time.Second {
		return true
	}

	// randomly sample 5% of remaining events
	return rand.Float64() < 0.05
}

var _ http.ResponseWriter = (*responseWrapper)(nil)
var _ http.Flusher = (*responseWrapper)(nil)
var _ http.Hijacker = (*responseWrapper)(nil)

type responseWrapper struct {
	delegate   http.ResponseWriter
	statusCode int
	headers    bool // have headers been written?
}

func (w *responseWrapper) Header() http.Header { return w.delegate.Header() }

func (w *responseWrapper) Write(p []byte) (int, error) {
	if !w.headers {
		w.headers = true
		w.statusCode = http.StatusOK
	}

	return w.delegate.Write(p)
}

func (w *responseWrapper) WriteHeader(statusCode int) {
	if !w.headers {
		w.headers = true
		w.statusCode = statusCode
		w.delegate.WriteHeader(statusCode)
	}
}

func (w *responseWrapper) Flush() {
	if f, ok := w.delegate.(http.Flusher); ok {
		if !w.headers {
			w.headers = true
			w.statusCode = http.StatusOK
		}

		f.Flush()
	}
}

func (w *responseWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := w.delegate.(http.Hijacker); ok {
		return h.Hijack()
	}

	return nil, nil, errors.New("hijack is not supported")
}
