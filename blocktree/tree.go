package blocktree

import "fmt"

// Tree
type Tree map[Hash]*Block

// Create new tree
func NewTree() (tree Tree, genesis Hash) {
	tree = Tree{}
	genesis = tree.NewBlock("genesis", NilHash)
	return
}

// Add block to tree
func (t *Tree) NewBlock(data string, previous Hash) Hash {
	block := Block{[]byte(data), previous, 0}
	block.Mine()
	hash := block.CalculateHash()
	(*t)[hash] = &block
	return hash
}

// View tree
func (t *Tree) View(branch Hash) {
	for b, i := branch, 0; b != NilHash; b, i = (*t)[b].Previous, i+1 {
		if (*t)[b].Previous == NilHash {
			fmt.Printf("#ROOT\t[%x]\n", b)
		} else {
			fmt.Printf("#%d\t[%x]\n", i, b)
		}
	}
}
