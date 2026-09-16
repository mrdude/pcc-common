package pcommon

import (
	"context"
	"os"
	"os/signal"
)

// WithSIGINTSubcontext creates a subcontext that cancels itself when SIGINT is received
// On Linux, SIGTERM is also handled.
func WithSIGINTSubcontext(ctx context.Context) context.Context {
	subCtx, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()

		sigCh := make(chan os.Signal, 1) // this channel must be buffered: https://stackoverflow.com/questions/68593779/can-unbuffered-channel-be-used-to-receive-signal/68593935#68593935
		signal.Notify(sigCh, interruptSignals...)

		for {
			select {
			case <-ctx.Done():
				if err := ctx.Err(); err != nil {
					return
				}
			case <-sigCh:
				return
			}
		}
	}()
	return subCtx
}
