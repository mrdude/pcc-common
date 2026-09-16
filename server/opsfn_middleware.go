package server

import (
	"github.com/gorilla/mux"
)

func WithGlobalMiddleware(mwf ...mux.MiddlewareFunc) OptionsFn {
	return OptionsFn(func(opt *Options) {
		if len(mwf) > 0 {
			opt.GlobalMiddlewares = append(opt.GlobalMiddlewares, mwf...)
		}
	})
}

func WithMiddleware(mwf ...mux.MiddlewareFunc) OptionsFn {
	return OptionsFn(func(opt *Options) {
		if len(mwf) > 0 {
			opt.Middlewares = append(opt.Middlewares, mwf...)
		}
	})
}
