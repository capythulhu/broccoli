package blocktree

import (
	"github.com/thzoid/broccoli/hash"
	"github.com/thzoid/broccoli/wallet"
)

// Transaction Output struct
type TxOutput struct {
	Address wallet.Address
	Value   uint64
}

// Transaction Input struct
type TxInput struct {
	ID    hash.Hash
	Index uint8
}

// Check if the address is the owner of the output
func (out *TxOutput) CanBeUnlocked(data wallet.Address) bool {
	return out.Address == data
}

// Transaction struct
type Transaction struct {
	Inputs  []TxInput
	Outputs [2]TxOutput
}

// Get transaction hash
func (tx Transaction) Hash() hash.Hash {
	buffer := []byte{}
	// Append inputs to buffer
	for _, in := range tx.Inputs {
		// Append input ID
		buffer = append(buffer, in.ID[:]...)
		// Append output ID
		buffer = append(buffer, toBytes(in.Index)...)
	}

	// Append outputs to buffer
	for _, out := range tx.Outputs {
		// Append output ID
		buffer = append(buffer, out.Address[:]...)
		// Append signature
		buffer = append(buffer, toBytes(out.Value)...)
	}

	// Build hash
	hash := hash.Hash{}
	hash.SHA256(buffer)
	return hash
}

// Check if transaction originated from coinbase
func (tx Transaction) IsFromCoinbase() bool {
	return len(tx.Inputs) == 1 && tx.Inputs[0].Index == 0 && tx.Hash() == hash.NilHash
}
