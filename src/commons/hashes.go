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

// Hashable is a technical interface to calculate a hash of a struct or element.
// It means we may
type Hashable interface {
	// ToHashString returns the hash of the entity.
	// Because an entity is immutable, the hash string should be invariant.
	ToHashString() string
}
