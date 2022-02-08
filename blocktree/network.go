package blocktree

import (
	"math/big"
)

// Network struct
type Network struct {
	Difficulty byte
	Reward     uint64
	UTXOs      map[string]TxOutput
}

// Build difficulty BigInt
func (n Network) BuildDifficultyBigInt() (target *big.Int) {
	target = big.NewInt(1)
	target.Lsh(target, 0xff-uint(n.Difficulty))
	return
}
