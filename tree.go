package main

import "fmt"

// Tree
type Tree map[Hash]Block

// Add block to tree
func (t *Tree) NewBlock(data string, previous Hash) Hash {
	block := Block{[]byte(data), previous}
	hash := block.CalculateHash()
	(*t)[hash] = block
	return hash
}

// View tree
func (t *Tree) View(branch Hash) {
	for b, i := branch, 0; b != NilHash; b, i = (*t)[b].Previous, i+1 {
		fmt.Printf("#%d [%x]\n", i, b)
	}
}
