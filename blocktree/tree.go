package blocktree

import (
	"encoding/hex"
	"fmt"
)

// Blocktree struct
type Blocktree struct {
	Blocks  map[Hash]*Block
	Network *Network
}

// Create new tree
func NewTree(n *Network, address string) (bt Blocktree, root Hash) {
	bt = Blocktree{map[Hash]*Block{}, n}
	coinbaseTx := GenerateRewardTx(n, address)
	root = bt.NewBlock([]*Transaction{coinbaseTx}, NilHash)
	return
}

// Add block to tree
func (bt *Blocktree) NewBlock(txs []*Transaction, previous Hash) Hash {
	block := Block{txs, previous, 0}
	block.Mine(bt.Network)
	hash := block.CalculateHash()
	bt.Blocks[hash] = &block
	return hash
}

// View tree
func (bt *Blocktree) View(branch Hash) {
	for i, b := 0, branch; b != NilHash; i, b = i+1, bt.Blocks[b].Previous {
		if bt.Blocks[b].Previous == NilHash {
			fmt.Printf("╳─ #r\t%x\n", b)
		} else if i == 0 {
			fmt.Printf("┌─ #%d\t%x\n", i, b)
		} else {
			fmt.Printf("├─ #%d\t%x\n", i, b)
		}
	}
}

// Get list of unspent transactions
func (bt *Blocktree) FindUnspentTxs(address string, branch Hash) []Transaction {
	var unspentTxs []Transaction
	spentTxNs := make(map[string][]int)

	for b := branch; b != NilHash; b = bt.Blocks[b].Previous {
		for _, tx := range bt.Blocks[b].Transactions {
			txID := hex.EncodeToString(tx.ID[:])
		Outputs:
			for i, out := range tx.Outputs {
				if spentTxNs[txID] != nil {
					for _, spentOut := range spentTxNs[txID] {
						if spentOut == i {
							continue Outputs
						}
					}
				}
				if out.CanBeUnlocked(address) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}
			if !tx.IsCoinbase() {
				for _, in := range tx.Inputs {
					if in.CanUnlock(address) {
						inTxID := hex.EncodeToString(in.ID[:])
						spentTxNs[inTxID] = append(spentTxNs[inTxID], in.OutID)
					}
				}
			}
		}
	}

	return unspentTxs
}

func (tree *Blocktree) FindUTxO(address string, branch Hash) []TxOutput {
	var UTXOs []TxOutput
	unspentTransactions := tree.FindUnspentTxs(address, branch)
	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

func (tree *Blocktree) FindSpendableOutputs(address string, amount uint64, branch Hash) (uint64, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := tree.FindUnspentTxs(address, branch)
	accumulated := uint64(0)

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString([]byte(tx.ID[:]))
		for i, out := range tx.Outputs {
			if out.CanBeUnlocked(address) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], i)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}
	return accumulated, unspentOuts
}
