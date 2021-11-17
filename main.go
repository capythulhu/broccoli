package main

import (
	"fmt"

	"github.com/thzoid/broccoli/blocktree"
)

func main() {
	tree, genesis := blocktree.NewTree()
	branchA1 := tree.NewBlock("foo block", genesis)
	branchA2 := tree.NewBlock("bar block", branchA1)
	branchA3 := tree.NewBlock("baz block", branchA2)
	branchB1 := tree.NewBlock("strange block", genesis)
	branchB2 := tree.NewBlock("odd block", branchB1)
	fmt.Println("\nBranch A")
	tree.View(branchA3)
	fmt.Println("\nBranch B")
	tree.View(branchB2)
	fmt.Println("is genesis valid?", tree[genesis].ValidateSimple())
}
