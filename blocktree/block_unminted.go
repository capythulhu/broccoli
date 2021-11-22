package blocktree

import (
	"errors"
)

// Unminted block struct
type UnmintedBlock struct {
	Transactions []Transaction
	Previous     Hash
}

// Create new block data
func NewBlock(previous Hash) UnmintedBlock {
	return UnmintedBlock{[]Transaction{}, previous}
}

// Add transaction
func (b *UnmintedBlock) AddTx(tree Blocktree, from, to string, amount uint64) error {
	var inputs []TxInput
	var outputs []TxOutput

	// Check if there are funds to perform transaction
	acc, validOutputs := tree.findSpendableOutputs(from, amount, b.Previous)
	if acc < amount {
		return errors.New("not enough funds")
	}

	// Generate inputs and outputs
	for tx, outs := range validOutputs {
		for _, out := range outs {
			input := TxInput{tx, out, from}
			inputs = append(inputs, input)
		}
	}

	// Append outputs
	outputs = append(outputs, TxOutput{amount, to})
	if acc > amount {
		outputs = append(outputs, TxOutput{acc - amount, from})
	}

	// Create and append transactions
	tx := Transaction{Inputs: inputs, Outputs: outputs}
	b.Transactions = append(b.Transactions, tx)
	return nil
}

// Add reward transaction
func (b *UnmintedBlock) AddRewardTx(tree Blocktree, to string) {
	// Coinbase input
	input := TxInput{NilHash, 0, "coinbase"}
	// Miner output
	output := TxOutput{tree.network.Reward, to}
	// Transaction
	tx := Transaction{Inputs: []TxInput{input}, Outputs: []TxOutput{output}}
	// Append transaction
	b.Transactions = append(b.Transactions, tx)
}
