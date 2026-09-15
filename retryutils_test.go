package pcommon

import (
	"testing"
	"time"
)

func TestIntPow2(t *testing.T) {
	i := 0
	expected := 1

	for ; i < 65; i++ {
		actual := intPow2(i)
		if expected != actual {
			t.Fatalf("intPow2(%d) is wrong (actual = %d, expected = %d)",
				i, actual, expected)
		}

		t.Logf("2^%d = %d", i, actual)

		expected *= 2
	}
}

func TestCalculateDelay(t *testing.T) {
	r := &Retryer[any]{
		Delay:    2 * time.Second,
		MaxDelay: 5 * time.Minute,

		noJitter: true,
	}

	for attempts := 0; attempts < 65; attempts++ {
		t.Logf("r.calculateDelay(attempts=%d) => %s",
			attempts, r.calculateDelay(attempts).String())
	}
}
