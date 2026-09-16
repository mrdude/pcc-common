//go:build !linux

package pcommon

import (
	"os"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
}
