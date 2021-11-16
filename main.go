package main

import "fmt"

func main() {
	tree := Tree{}
	genesis := tree.NewBlock("genesis block", NilHash)
	branchA1 := tree.NewBlock("foo block", genesis)
	branchA2 := tree.NewBlock("bar block", branchA1)
	branchA3 := tree.NewBlock("baz block", branchA2)
	branchB1 := tree.NewBlock("strange block", genesis)
	branchB2 := tree.NewBlock("odd block", branchB1)
	fmt.Println("Branch A")
	tree.View(branchA3)
	fmt.Println("Branch B")
	tree.View(branchB2)
}
