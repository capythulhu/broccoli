package blocktree

import "fmt"

// Blocktree struct
type Blocktree struct {
	Blocks  map[Hash]*Block
	Network *Network
}

// Create new tree
func NewTree(n *Network) (bt Blocktree, genesis Hash) {
	bt = Blocktree{map[Hash]*Block{}, n}
	genesis = bt.NewBlock("genesis", NilHash)
	return
}

// Add block to tree
func (bt *Blocktree) NewBlock(data string, previous Hash) Hash {
	block := Block{[]byte(data), previous, 0}
	block.Mine(bt.Network)
	hash := block.CalculateHash()
	bt.Blocks[hash] = &block
	return hash
}

// View tree
func (bt *Blocktree) View(branch Hash) {
	for b, i := branch, 0; b != NilHash; b, i = bt.Blocks[b].Previous, i+1 {
		if bt.Blocks[b].Previous == NilHash {
			fmt.Printf("#root\t%x\n", b)
		} else {
			fmt.Printf("#%d\t%x\n", i, b)
		}
	}
}
