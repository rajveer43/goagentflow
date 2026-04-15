package idempotency

import (
	"crypto/rand"
	"encoding/hex"
)

func NewKey() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
