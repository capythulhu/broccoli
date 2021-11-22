package blocktree

// Transaction Output struct
type TxOutput struct {
	Value  uint64
	PubKey string
}

// Transaction Input struct
type TxInput struct {
	ID        Hash
	OutID     int
	Signature string
}

// Check if the address is the owner of the input
func (in *TxInput) CanUnlock(data string) bool {
	return in.Signature == data
}

// Check if the address is the owner of the output
func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data
}
