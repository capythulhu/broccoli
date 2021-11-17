package blocktree

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"math"
	"math/big"
)

// Hash
type Hash [32]byte

// Nil Hash
var NilHash = Hash{}

// Block
type Block struct {
	Data     []byte
	Previous Hash
	Nonce    uint
}

// Int to bytes
func toBytes(nonce uint) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, int64(nonce))
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

// Calculate block hash
func (b *Block) CalculateHash() Hash {
	info := bytes.Join([][]byte{b.Data, b.Previous[:], toBytes(b.Nonce)}, []byte{})
	return sha256.Sum256(info)
}

// Mine block
func (b *Block) Mine(n *Network) {
	target := n.BuildDifficultyBigInt()
	intHash := big.NewInt(0)
	b.Nonce = 0
	for b.Nonce < math.MaxInt64 {
		if b.Validate(target, intHash) {
			break
		} else {
			b.Nonce++
		}

	}
}

// Validate block nonce
func (b *Block) Validate(target *big.Int, intHash *big.Int) bool {
	hash := b.CalculateHash()
	intHash.SetBytes(hash[:])

	return intHash.Cmp(target) == -1
}

// Validate block nonce generating buffers
func (b *Block) ValidateSimple(n *Network) bool {
	return b.Validate(n.BuildDifficultyBigInt(), big.NewInt(0))
}
