package server

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type OptionsFn func(*Options)

type Options struct {
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration

	MaxHeaderBytes int

	TLSConfig *tls.Config

	GlobalMiddlewares []mux.MiddlewareFunc //GlobalMiddlewares get executed even for NotFoundHandlers
	Middlewares       []mux.MiddlewareFunc
}

func getDefaultOptions() Options {
	return Options{
		ReadTimeout:       0,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      0,
		IdleTimeout:       5 * time.Minute,
		MaxHeaderBytes:    1024 * 1024 * 4,

		TLSConfig: nil,

		GlobalMiddlewares: nil,
		Middlewares:       nil,
	}
}

func applyOptions(fn ...OptionsFn) Options {
	op := getDefaultOptions()
	for _, f := range fn {
		f(&op)
	}

	return op
}

func (op *Options) createHttpServer() *http.Server {
	return &http.Server{
		ReadTimeout:       op.ReadTimeout,
		ReadHeaderTimeout: op.ReadHeaderTimeout,
		WriteTimeout:      op.WriteTimeout,
		IdleTimeout:       op.IdleTimeout,
		MaxHeaderBytes:    op.MaxHeaderBytes,
		TLSConfig:         op.TLSConfig,
	}
}

func (op *Options) addRouterMiddleware(r *mux.Router) {
	for _, m := range op.Middlewares {
		r.Use(m)
	}
}

func wrapHandlerWithMiddleware(h http.Handler, mw []mux.MiddlewareFunc) http.Handler {
	wrappedHandler := h
	for i := len(mw) - 1; i >= 0; i-- {
		mwf := mw[i]
		wrappedHandler = mwf.Middleware(wrappedHandler)
	}

	return wrappedHandler
}
