package server

import (
	"context"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mrdude/pcc-common"
	"github.com/mrdude/pcc-common/server/serverctx"
)

type httpServerImpl struct {
	name    string
	sv      *http.Server
	handler http.Handler

	middleware       []mux.MiddlewareFunc
	globalMiddleware []mux.MiddlewareFunc
}

func CreateHttpServer(name string, handler http.Handler, options ...OptionsFn) HttpServer {
	op := applyOptions(options...)

	sv := &httpServerImpl{
		name:    name,
		sv:      op.createHttpServer(),
		handler: handler,
	}

	sv.middleware = append(sv.middleware, op.Middlewares...)
	sv.globalMiddleware = append(sv.globalMiddleware, op.GlobalMiddlewares...)

	return sv
}

func (sv *httpServerImpl) Name() string {
	return sv.name
}

func (sv *httpServerImpl) HttpServer() *http.Server {
	return sv.sv
}

func (sv *httpServerImpl) Run(ctx context.Context, sock net.Listener) error {
	//attach this server's name to the logger
	logger := pcommon.GetLogger(ctx)
	logger = attachServerNameToLogger(logger, sv)
	ctx = pcommon.WithLoggerContextValue(ctx, logger)

	//wrap the handler with middleware
	wrappedHandler := sv.handler
	wrappedHandler = wrapHandlerWithMiddleware(wrappedHandler, sv.middleware)
	wrappedHandler = wrapHandlerWithMiddleware(wrappedHandler, sv.globalMiddleware)

	//initialize the http server
	sv.sv.BaseContext = createBaseContextFn(ctx)
	sv.sv.ConnContext = func(ctx context.Context, c net.Conn) context.Context {
		return serverctx.WithConnectionContextValue(ctx, c)
	}
	sv.sv.Handler = wrappedHandler

	if sv.sv.TLSConfig == nil {
		return sv.sv.Serve(sock)
	} else {
		return sv.sv.ServeTLS(sock, "", "")
	}
}

func (sv *httpServerImpl) Shutdown(ctx context.Context) error {
	return sv.sv.Shutdown(ctx)
}
