package blocktree

// Transaction Output struct
type TxOutput struct {
	PubKey string
	Value  uint64
}

// Transaction Input struct
type TxInput struct {
	ID    Hash
	Index uint8
}

// Check if the address is the owner of the output
func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data
}
