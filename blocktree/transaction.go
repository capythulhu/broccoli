package blocktree

import (
	"crypto/sha256"
)

// Transaction struct
type Transaction struct {
	Inputs  []TxInput
	Outputs []TxOutput
}

// Get transaction hash
func (tx *Transaction) Hash() Hash {
	buffer := []byte{}
	// Append inputs to buffer
	for _, in := range tx.Inputs {
		// Append input ID
		buffer = append(buffer, in.ID[:]...)
		// Append output ID
		buffer = append(buffer, toBytes(int64(in.Index))...)
	}

	// Append outputs to buffer
	for _, out := range tx.Outputs {
		// Append output ID
		buffer = append(buffer, out.PubKey...)
		// Append signature
		buffer = append(buffer, toBytes(int64(out.Value))...)
	}

	return DecodeSHA256(sha256.Sum256(buffer))
}

// Check if transaction originated from coinbase
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && tx.Inputs[0].Index == 0 && tx.Hash() == NilHash
}
