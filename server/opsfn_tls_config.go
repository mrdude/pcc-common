package server

import (
	"crypto/tls"
)

func WithTLSConfig(cfg *tls.Config) OptionsFn {
	return OptionsFn(func(opt *Options) {
		opt.TLSConfig = cfg
	})
}
