package main

import (
	"fmt"

	"github.com/thzoid/broccoli/blocktree"
)

func BlockView(bt *blocktree.Blocktree, h blocktree.Hash) {
	block := bt.FindBlock(h)
	if block.Previous() == blocktree.NilHash {
		fmt.Printf("block %x (root)\n", h)
	} else {
		fmt.Printf("block %x\n", h)
	}
	for sender, tx := range block.Transactions() {
		fmt.Printf("├─tx %x\t\t ◀ %s\n", tx.Hash(), sender)
		txConnector := '├'
		for j, in := range tx.Inputs {
			if j == len(tx.Inputs)-1 && len(tx.Outputs) == 0 {
				txConnector = '└'
			}
			fmt.Printf("│ %c─in  %x[%d]\n", txConnector, in.ID, in.Index)
		}
		noChange := tx.Outputs[1].Value == 0
		for j, out := range tx.Outputs {
			if j == 0 && noChange || j == 1 {
				txConnector = '└'
			}
			if out.Value > 0 {
				fmt.Printf("│ %c─out %d\t\t ▶ %s\n", txConnector, out.Value, out.PubKey)
			}
		}
		fmt.Printf("│\n")
	}
}

func main() {
	tree, root := blocktree.NewTree(blocktree.Network{Difficulty: 8, Reward: 100}, "alice")
	BlockView(&tree, root)

	b1 := blocktree.NewBlock(root)
	b1.AddTx(tree, "alice", blocktree.TxOutput{PubKey: "bob", Value: 5})
	b1Hash := tree.Graft(b1, "carol")

	BlockView(&tree, b1Hash)

	// b2 := blocktree.NewBlock(b1Hash)
	// b2.AddTx(tree, "bob", blocktree.TxOutput{PubKey: "dave", Value: 3})
	// b2.AddTx(tree, "carol", blocktree.TxOutput{PubKey: "dave", Value: 2})
	// b2Hash := tree.Graft(b2, "carol")
	// BlockView(&tree, b2Hash)

	// b3 := blocktree.NewBlock(b2Hash)
	// b3.AddTx(tree, "carol", blocktree.TxOutput{PubKey: "alice", Value: 198})
	// b3Hash := tree.Graft(b3, "carol")
	// BlockView(&tree, b3Hash)

	// fmt.Println("\nTree view")
	// tree.View(b3Hash)
}
