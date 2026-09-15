package pcommon

import (
	"context"
	"errors"
	"math/rand"
	"time"
)

// Retryer handles exponential retry with jitter
type Retryer[T any] struct {
	Delay       time.Duration                        // the base amount of delay between calls
	MaxDelay    time.Duration                        // the maximum delay. set to 0 to ensure no max delay (unrecommended!)
	MaxAttempts int                                  // the max number of attempts to make. set to 0 for no maximum attempts
	Fn          func(ctx context.Context) (T, error) // the callback function

	noJitter bool // when set to true, we won't apply jitter. this is used for testing
}

// Do executes the Retryer.Fn callback with exponential retry and jitter.
//
// Errors of type context.Canceled or context.DeadlineExceeded will be returned immediately.
// All other errors will be suppressed until we've reached max number of attempts.
func (r Retryer[T]) Do(ctx context.Context) (T, error) {
	// verify parameters
	if r.Delay <= 0 {
		var t T
		return t, errors.New(".Delay must be > 0")
	}

	if r.Fn == nil {
		var t T
		return t, errors.New(".Fn must be set")
	}

	// main loop
	attempts := 0 // number of consecutive failed attempts
	for {
		if err := ctx.Err(); err != nil {
			// Our context has expired -- this error won't go away with retries
			var t T
			return t, err
		}

		ret, err := r.Fn(ctx)
		if err == nil {
			return ret, nil
		} else if r.MaxAttempts != 0 && attempts == r.MaxAttempts {
			// we've hit the limit on failures -- bail out
			var t T
			return t, err
		} else {
			// swallow the error and continue

			if err = r.nonBlockingSleep(ctx, r.calculateDelay(attempts)); err != nil {
				var t T
				return t, err
			}

			attempts++
		}
	}
}

func (r Retryer[T]) calculateDelay(attempts int) time.Duration {
	if attempts < 0 {
		panic(errors.New("attempts must be >0"))
	}

	delay := r.Delay * time.Duration(intPow2(attempts))
	if delay <= 0 {
		// handle wraparound

		const maxDuration time.Duration = 1<<63 - 1 // taken from stdlib's time.go
		delay = maxDuration
	}

	// enforce max delay
	if r.MaxDelay != 0 && delay > r.MaxDelay {
		delay = r.MaxDelay
	}

	// apply jitter
	if !r.noJitter {
		delay = randDuration(delay)
	}

	return delay
}

// TODO make this an exported function?
func randDuration(N time.Duration) time.Duration {
	X := rand.Int63n(int64(N))
	return time.Duration(X)
}

func (r Retryer[T]) nonBlockingSleep(ctx context.Context, d time.Duration) error {
	t := time.NewTicker(d)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				return err
			}
		case <-t.C:
			return nil
		}
	}
}

// returns math.Pow(2, exp)
func intPow2(exp int) int {
	switch {
	case exp < 0:
		return 0
	case exp == 0:
		return 1
	case exp == 1:
		return 2
	case exp%2 == 0: // even
		v := intPow2(exp / 2)
		return v * v
	default: // exp%2 != 0: // odd
		v := intPow2((exp - 1) / 2)
		return 2 * v * v
	}
}
