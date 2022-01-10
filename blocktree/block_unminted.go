package blocktree

import (
	"errors"

	"github.com/elliotchance/orderedmap"
	"github.com/thzoid/broccoli/hash"
	"github.com/thzoid/broccoli/wallet"
)

// Unminted block struct
type UnmintedBlock struct {
	Transactions orderedmap.OrderedMap
	Previous     hash.Hash
}

// Create new block data
func NewBlock(previous hash.Hash) UnmintedBlock {
	return UnmintedBlock{*orderedmap.NewOrderedMap(), previous}
}

// Add transaction
func (b *UnmintedBlock) AddTx(tree Blocktree, from wallet.Address, output TxOutput) error {
	var inputs []TxInput

	// Calculate total amount
	amount := output.Value

	// Check if there are funds to perform transaction
	acc, validOutputs := tree.findSpendableOutputs(from, amount, b.Previous)
	if acc < amount {
		return errors.New("not enough funds")
	}

	// Generate inputs and outputs
	for tx, outs := range validOutputs {
		for _, out := range outs {
			input := TxInput{tx, out}
			inputs = append(inputs, input)
		}
	}

	// Transaction final outputs
	outputs := [2]TxOutput{
		output,               // Receipient
		{from, acc - amount}, // Change
	}

	// Create and append transactions
	tx := Transaction{Inputs: inputs, Outputs: outputs}
	b.Transactions.Set(from, tx)
	return nil
}

// Add reward transaction
func (b *UnmintedBlock) AddRewardTx(tree Blocktree, to wallet.Address) {
	// Coinbase input
	input := TxInput{hash.NilHash, 0}
	// Miner output
	output := TxOutput{to, tree.network.Reward}
	// Transaction
	tx := Transaction{Inputs: []TxInput{input}, Outputs: [2]TxOutput{output}}
	// Append transaction
	b.Transactions.Set(wallet.CoinbaseAddress(), tx)
}
