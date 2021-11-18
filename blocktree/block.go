package blocktree

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"log"
	"math"
	"math/big"
)

// Block
type Block struct {
	Transactions []*Transaction
	Previous     Hash
	Nonce        uint
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

// Serialize block
func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	panic(err)
	return res.Bytes()
}

// Get Transactions hash
func (b *Block) CalculateTxsHash() Hash {
	var hashes [][]byte
	for _, t := range b.Transactions {
		hashes = append(hashes, t.ID[:])
	}
	return sha256.Sum256(bytes.Join(hashes, []byte{}))
}

// Calculate block hash
func (b *Block) CalculateHash() Hash {
	txsHash := b.CalculateTxsHash()
	info := bytes.Join([][]byte{b.Previous[:], toBytes(b.Nonce), txsHash[:]}, []byte{})
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
