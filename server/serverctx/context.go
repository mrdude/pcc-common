package serverctx

import (
	"context"
	"net"
)

type contextKey string

const (
	connContextKey     contextKey = "CONNECTION"
	listenerContextKey contextKey = "LISTENER"
)

func WithConnectionContextValue(ctx context.Context, conn net.Conn) context.Context {
	return context.WithValue(ctx, connContextKey, conn)
}

func GetConnection(ctx context.Context) net.Conn {
	v := ctx.Value(connContextKey)
	if v == nil {
		return nil
	}
	return v.(net.Conn)
}

func WithListenerContextValue(ctx context.Context, L net.Listener) context.Context {
	return context.WithValue(ctx, listenerContextKey, L)
}

func GetListener(ctx context.Context) net.Listener {
	v := ctx.Value(listenerContextKey)
	if v == nil {
		return nil
	}
	return v.(net.Listener)
}
