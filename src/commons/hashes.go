package commons

import (
	"crypto/sha512"
	"encoding/hex"
)

// HashString returns a SHA-512 hash of the given string.
// Key points : same text => same hash and different text => different hash in practice
func HashString(text string) string {
	hasher := sha512.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
