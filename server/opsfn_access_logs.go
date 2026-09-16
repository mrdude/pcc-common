package server

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mrdude/pcc-common"
	"go.uber.org/zap"
)

// WithAccessLogs adds a global middleware that logs every HTTP request and response
func WithAccessLogs() OptionsFn {
	mwf := mux.MiddlewareFunc(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			//wrap the response writer and start the timer
			wrappedResp := wrapResponseWriter(w)
			startTime := time.Now()

			//record the response
			logger := pcommon.GetLogger(req.Context())
			defer func() {
				endTime := time.Now()
				requestDuration := endTime.Sub(startTime)

				panicErr := recover()

				fields := []zap.Field{
					zap.Bool("access_log", true),
					zap.Duration("job.duration", requestDuration),
					zap.String("req.method", req.Method),
					zap.String("req.url", req.URL.String()),
					zap.String("req.proto", req.Proto),
				}

				if panicErr == nil {
					fields = append(fields,
						zap.Int("resp.status", wrappedResp.statusCode),
						zap.Int64("resp.body_len", wrappedResp.payloadBytesWritten),
					)
				} else {
					fields = append(fields,
						zap.Any("panic", panicErr),
					)
				}

				logger.Info("completed request", fields...)

				//repanic after logging the request
				if panicErr != nil {
					panic(panicErr)
				}
			}()

			//invoke the next handler
			next.ServeHTTP(wrappedResp, req)
		})
	})

	return OptionsFn(func(opt *Options) {
		opt.GlobalMiddlewares = append(opt.GlobalMiddlewares, mwf)
	})
}

var _ http.ResponseWriter = (*responseWriterWrapper)(nil)
var _ http.Flusher = (*responseWriterWrapper)(nil)

type responseWriterWrapper struct {
	delegate http.ResponseWriter

	statusCode          int
	wroteHeader         bool
	payloadBytesWritten int64 //the number of non-header bytes written
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriterWrapper {
	return &responseWriterWrapper{
		delegate:            w,
		statusCode:          http.StatusOK,
		wroteHeader:         false,
		payloadBytesWritten: 0,
	}
}

func (w *responseWriterWrapper) doWriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.delegate.WriteHeader(statusCode)
}

func (w *responseWriterWrapper) Header() http.Header {
	return w.delegate.Header()
}

func (w *responseWriterWrapper) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.wroteHeader = true
		w.doWriteHeader(http.StatusOK)
	}

	n, err := w.delegate.Write(b)
	w.payloadBytesWritten += int64(n)
	return n, err
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	if !w.wroteHeader {
		w.wroteHeader = true
		w.doWriteHeader(statusCode)
	}
}

func (w *responseWriterWrapper) Flush() {
	if f, ok := w.delegate.(http.Flusher); ok {
		if !w.wroteHeader {
			w.wroteHeader = true
			w.doWriteHeader(http.StatusOK)
		}

		f.Flush()
	}
}

func (w *responseWriterWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := w.delegate.(http.Hijacker); ok {
		return h.Hijack()
	}

	return nil, nil, errors.New("hijack is not supported")
}

func (w *responseWriterWrapper) SentStatusCode() int {
	return w.statusCode
}

func (w *responseWriterWrapper) SentPayloadBytes() int64 {
	return w.payloadBytesWritten
}
