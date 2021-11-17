package blocktree

import "math/big"

// Network struct
type Network struct {
	Difficulty byte
}

// Build difficulty BigInt
func (n *Network) BuildDifficultyBigInt() (target *big.Int) {
	target = big.NewInt(1)
	target.Lsh(target, uint(255-n.Difficulty))
	return
}
