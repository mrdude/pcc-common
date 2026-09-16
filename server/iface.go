package server

import (
	"context"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

type Server interface {
	Name() string
	Run(ctx context.Context, sock net.Listener) error
}

// HttpServer is a thin wrapper around http.Server
type HttpServer interface {
	Server
	HttpServer() *http.Server
}

// WebServer is a wrapper around HttpServer that provides a mux.Router as a handler
type WebServer interface {
	HttpServer

	Router() *mux.Router
	Shutdown(ctx context.Context) error
}
