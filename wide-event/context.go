package wideevt

import "context"

type contextKey string

const eventContextKey contextKey = "event"

func WithEvent(ctx context.Context, evt *Event) context.Context {
	return context.WithValue(ctx, eventContextKey, evt)
}

func GetEvent(ctx context.Context) *Event {
	v := ctx.Value(eventContextKey)
	if v == nil {
		return nil
	}
	return v.(*Event)
}
