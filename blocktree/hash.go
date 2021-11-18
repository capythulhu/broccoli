package blocktree

// Hash
type Hash [32]byte

// Nil Hash
var NilHash = Hash{}

// Fill Hash with bytes slice
func HashFromBytes(b []byte) (h Hash) {
	copy(h[:], b[:32])
	return
}
