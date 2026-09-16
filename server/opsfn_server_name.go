package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

// WithServerName sets the Server header on all outgoing requests.
func WithServerName(name string) OptionsFn {
	mwf := mux.MiddlewareFunc(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Server", name)

			//invoke the next handler
			next.ServeHTTP(w, req)
		})
	})

	return OptionsFn(func(opt *Options) {
		opt.GlobalMiddlewares = append(opt.GlobalMiddlewares, mwf)
	})
}
