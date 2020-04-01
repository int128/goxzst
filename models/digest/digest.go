// Package digest provides the message digest features.
package digest

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
)

var (
	SHA256 = &Algorithm{
		Name:   "sha256",
		Suffix: ".sha256",
		NewHash: func() hash.Hash {
			return sha256.New()
		},
	}
	SHA512 = &Algorithm{
		Name:   "sha512",
		Suffix: ".sha512",
		NewHash: func() hash.Hash {
			return sha512.New()
		},
	}
)

var AvailableAlgorithms = []*Algorithm{
	SHA256,
	SHA512,
}

// NewAlgorithm returns the corresponding Algorithm.
// If not found, it returns an error.
func NewAlgorithm(name string) (*Algorithm, error) {
	for _, alg := range AvailableAlgorithms {
		if name == alg.Name {
			return alg, nil
		}
	}
	return nil, fmt.Errorf("not supported %s", name)
}

// Algorithm represents a digest algorithm.
type Algorithm struct {
	Name    string
	Suffix  string
	NewHash func() hash.Hash
}
