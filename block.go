package main

import (
	"bytes"
	"crypto/sha256"
)

// Hash
type Hash [32]byte

// Nil Hash
var NilHash = Hash{}

// Block
type Block struct {
	Data     []byte
	Previous Hash
}

// Calculate block hash
func (b *Block) CalculateHash() [32]byte {
	info := bytes.Join([][]byte{b.Data, b.Previous[:]}, []byte{})
	return sha256.Sum256(info)
}
