package blocktree

import (
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"
)

// Minted block struct
type MintedBlock struct {
	tree *Blocktree
	// Actual block data
	transactions []Transaction
	previous     Hash
	nonce        uint
}

// Get previous block hash
func (b *MintedBlock) Previous() Hash {
	return b.previous
}

// Get block transactions
func (b *MintedBlock) Transactions() []Transaction {
	return b.transactions
}

// Calculate block hash
func (b *MintedBlock) Hash() Hash {
	// Calculate transactions hash
	var hashes [][]byte
	for _, t := range b.transactions {
		hash := t.Hash()
		hashes = append(hashes, hash[:])
	}
	txsHash := sha256.Sum256(bytes.Join(hashes, []byte{}))
	// Join block data
	data := bytes.Join([][]byte{b.previous[:], toBytes(int64(b.nonce)), txsHash[:]}, []byte{})
	return sha256.Sum256(data)
}

// Mine block
func (b *MintedBlock) mine(n Network) {
	target := n.BuildDifficultyBigInt()
	intHash := big.NewInt(0)
	b.nonce = 0
	for b.nonce < math.MaxInt64 {
		if b.validate(target, intHash) {
			break
		} else {
			b.nonce++
		}
	}
}

// Validate block nonce with external buffers
func (b *MintedBlock) validate(target *big.Int, intHash *big.Int) bool {
	hash := b.Hash()
	intHash.SetBytes(hash[:])

	return intHash.Cmp(target) == -1
}
