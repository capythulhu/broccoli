package main

import (
	"fmt"

	"github.com/thzoid/broccoli/blocktree"
)

func BlockView(bt *blocktree.Blocktree, h blocktree.Hash) {
	fmt.Printf("┌block %x\n", h)
	for i, tx := range bt.FindBlock(h).Transactions() {
		fmt.Printf("├─tx #%d id: %x\n", i+1, tx.Hash())
		for _, in := range tx.Inputs {
			fmt.Printf("│ ├─in %x[%d] (%s)\n", in.ID, in.OutID, in.Signature)
		}
		for _, out := range tx.Outputs {
			fmt.Printf("│ ├─out %d -> %s\n", out.Value, out.PubKey)
		}
	}
}

func main() {
	// cmd.Execute()
	tree, root := blocktree.NewTree(blocktree.Network{Difficulty: 16, Reward: 100}, "alice")

	b1a := blocktree.NewBlock(root)
	b1a.AddTx(tree, "alice", "bob", 80)
	b1aHash := tree.Graft(b1a, "eve")

	b1b := blocktree.NewBlock(root)
	b1b.AddTx(tree, "alice", "eve", 60)
	b1bHash := tree.Graft(b1b, "eve")

	BlockView(&tree, root)
	BlockView(&tree, b1aHash)
	BlockView(&tree, b1bHash)

	// fmt.Println("\nTree view")
	// tree.View(b1aHash)
}
