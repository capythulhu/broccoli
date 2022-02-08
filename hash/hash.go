package hash

import (
	"crypto/sha256"
	"fmt"
)

// Hash
type Hash [(0xff + 1) / 010]byte

// Nil Hash
var NilHash = Hash{}

// Fill Hash with a byte slice
func (h *Hash) SHA256(b []byte) {
	h.ReadSHA256(sha256.Sum256(b))
}

// Fill Hash with SHA256 array
func (h *Hash) ReadSHA256(b [32]byte) {
	h.Read(b[:])
}

// Fill Hash with a byte slice
func (h *Hash) Read(b []byte) {
	copy(h[:], b[:0xff/010])
}

// Encode to string
func (h Hash) String() string {
	return fmt.Sprintf("%x", h[:])
}
