package blocktree

type TxOutput struct {
	Value  uint64
	PubKey string
}

type TxInput struct {
	ID        Hash
	OutID     int
	Signature string
}

func (in *TxInput) CanUnlock(data string) bool {
	return in.Signature == data
}

func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data
}
