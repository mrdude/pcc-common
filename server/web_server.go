package server

import (
	"context"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mrdude/pcc-common"
	"github.com/mrdude/pcc-common/server/serverctx"
)

type webServerImpl struct {
	name string
	sv   *http.Server
	r    *mux.Router

	globalMiddleware []mux.MiddlewareFunc

	handler http.Handler
}

func CreateWebServer(name string, options ...OptionsFn) WebServer {
	op := applyOptions(options...)

	r := mux.NewRouter()
	sv := &webServerImpl{
		name: name,
		sv:   op.createHttpServer(),
		r:    r,
	}

	op.addRouterMiddleware(sv.r)
	sv.globalMiddleware = append(sv.globalMiddleware, op.GlobalMiddlewares...)

	return sv
}

func (sv *webServerImpl) Name() string {
	return sv.name
}

func (sv *webServerImpl) HttpServer() *http.Server {
	return sv.sv
}

func (sv *webServerImpl) Router() *mux.Router {
	return sv.r
}

func (sv *webServerImpl) Run(ctx context.Context, sock net.Listener) error {
	//attach this server's name to the logger
	logger := pcommon.GetLogger(ctx)
	logger = attachServerNameToLogger(logger, sv)
	ctx = pcommon.WithLoggerContextValue(ctx, logger)

	//wrap the handler with global middleware
	sv.handler = wrapHandlerWithMiddleware(sv.r, sv.globalMiddleware)

	//initialize the http server
	sv.sv.BaseContext = createBaseContextFn(ctx)
	sv.sv.ConnContext = func(ctx context.Context, c net.Conn) context.Context {
		return serverctx.WithConnectionContextValue(ctx, c)
	}
	sv.sv.Handler = sv.handler

	if sv.sv.TLSConfig == nil {
		return sv.sv.Serve(sock)
	} else {
		return sv.sv.ServeTLS(sock, "", "")
	}
}

func createBaseContextFn(baseCtx context.Context) func(listener net.Listener) context.Context {
	return func(listener net.Listener) context.Context {
		ctx := baseCtx
		ctx = serverctx.WithListenerContextValue(ctx, listener)
		return ctx
	}
}

func (sv *webServerImpl) Shutdown(ctx context.Context) error {
	return sv.sv.Shutdown(ctx)
}
