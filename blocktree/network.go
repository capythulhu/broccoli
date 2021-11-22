package blocktree

import "math/big"

// Network struct
type Network struct {
	Difficulty byte
	Reward     uint64
	UTXOs      map[string]TxOutput
}

// Build difficulty BigInt
func (n *Network) BuildDifficultyBigInt() (target *big.Int) {
	target = big.NewInt(1)
	target.Lsh(target, uint(HASH_BYTES*8-1-n.Difficulty))
	return
}
