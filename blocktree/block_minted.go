package blocktree

import (
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"

	"github.com/elliotchance/orderedmap"
	"github.com/thzoid/broccoli/hash"
	"github.com/thzoid/broccoli/wallet"
)

// Minted block struct
type MintedBlock struct {
	tree *Blocktree
	// Actual block data
	transactions orderedmap.OrderedMap
	previous     hash.Hash
	nonce        uint32
}

// Get previous block hash
func (b *MintedBlock) Previous() hash.Hash {
	return b.previous
}

// Get block transactions
func (b *MintedBlock) Transactions() map[wallet.Address]Transaction {
	// Copy map values
	result := map[wallet.Address]Transaction{}
	for el := b.transactions.Front(); el != nil; el = el.Next() {
		result[el.Key.(wallet.Address)] = el.Value.(Transaction)
	}
	return result
}

// Calculate block hash
func (b *MintedBlock) Hash() hash.Hash {
	// Calculate transactions hash
	hashes := make([][]byte, b.transactions.Len())
	for i, el := 0, b.transactions.Front(); el != nil; i, el = i+1, el.Next() {
		hash := el.Value.(Transaction).Hash()
		hashes[i] = hash[:]
	}
	txsHash := sha256.Sum256(bytes.Join(hashes, []byte{}))
	// Join block data
	data := bytes.Join([][]byte{b.previous[:], toBytes(b.nonce), txsHash[:]}, []byte{})
	// Build hash
	hash := hash.Hash{}
	hash.SHA256(data)
	return hash
}

// Mine block
func (b *MintedBlock) mine(n Network) {
	target := n.BuildDifficultyBigInt()
	intHash := big.NewInt(0)
	b.nonce = 0
	for b.nonce < math.MaxUint32 {
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
