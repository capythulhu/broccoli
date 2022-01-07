package blocktree

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"

	"github.com/elliotchance/orderedmap"
)

// Minted block struct
type MintedBlock struct {
	tree *Blocktree
	// Actual block data
	transactions orderedmap.OrderedMap
	previous     Hash
	nonce        uint32
}

// Get previous block hash
func (b *MintedBlock) Previous() Hash {
	return b.previous
}

// Get block transactions
func (b *MintedBlock) Transactions() map[string]Transaction {
	// Copy map values
	result := map[string]Transaction{}
	for el := b.transactions.Front(); el != nil; el = el.Next() {
		result[el.Key.(string)] = el.Value.(Transaction)
	}
	return result
}

// Calculate block hash
func (b *MintedBlock) Hash() Hash {
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
	hash := Hash{}
	hash.Read(sha256.Sum256(data))
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
	fmt.Println("nonce:", b.nonce)
}

// Validate block nonce with external buffers
func (b *MintedBlock) validate(target *big.Int, intHash *big.Int) bool {
	hash := b.Hash()
	intHash.SetBytes(hash[:])

	return intHash.Cmp(target) == -1
}
