package blocktree

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
)

// Transaction struct
type Transaction struct {
	ID      Hash
	Inputs  []TxInput
	Outputs []TxOutput
}

// Set ID of transaction
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	hash := NilHash

	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	panic(err)

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash
}

// Miner reward transaction
func GenerateRewardTx(n *Network, toAddress string) *Transaction {
	txIn := TxInput{NilHash, 0, ""}
	txOut := TxOutput{n.Reward, toAddress}

	tx := Transaction{NilHash, []TxInput{txIn}, []TxOutput{txOut}}

	return &tx
}

// Check if transaction originated from coinbase
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && tx.Inputs[0].OutID == 0 && tx.ID == NilHash
}

// Create a new transaction
func NewTx(from, to string, amount uint64, branch Hash, tree *Blocktree) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	acc, validOutputs := tree.FindSpendableOutputs(from, amount, branch)

	if acc < amount {
		log.Panic("error: not enough funds!")
	}
	for tx, outs := range validOutputs {
		txIDSlice, err := hex.DecodeString(tx)
		if err != nil {
			panic(err)
		}

		txID := HashFromBytes(txIDSlice)
		for _, out := range outs {
			input := TxInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TxOutput{amount, to})

	if acc > amount {
		outputs = append(outputs, TxOutput{acc - amount, from})
	}

	tx := Transaction{NilHash, inputs, outputs}
	tx.SetID()

	return &tx
}
