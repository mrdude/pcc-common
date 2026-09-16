package pcommon

import (
	crand "crypto/rand"
	"errors"
)

func GenerateCryptoRandomBytes(n int) []byte {
	if n < 0 {
		panic(errors.New("invalid value"))
	}

	b := make([]byte, n)
	_, err := crand.Read(b)
	if err != nil {
		panic(err)
	}

	return b
}
