package blocktree

import (
	"errors"
	"fmt"

	"github.com/elliotchance/orderedmap"
)

// Blocktree struct
// This structure exists only to ensure
// the integrity of the tree (i.e. not
// allowing other functions to modify
// the actual blocks)
type Blocktree struct {
	blocks  map[Hash]MintedBlock
	network Network
}

// Create new tree
func NewTree(n Network, miner string) (bt Blocktree, root Hash) {
	bt = Blocktree{map[Hash]MintedBlock{}, n}
	// Create root block
	root = bt.Graft(NewBlock(NilHash), miner)
	return
}

// Perform validations and return expected minted block to be grafted
func (bt *Blocktree) Mint(ub UnmintedBlock, miner string) (*MintedBlock, error) {
	// Check if block is root but there are blocks on the tree already
	if ub.Previous == NilHash {
		if len(bt.blocks) != 0 {
			return nil, errors.New("multiple root blocks are not allowed")
		}
		// Check if block exists
	} else if bt.FindBlock(ub.Previous) == nil {
		return nil, errors.New("previous block not found")
	}
	// Check transactions
	for _, t := range ub.Transactions {
		if t.IsCoinbase() {
			return nil, errors.New("illegal coinbase transaction found")
		}
	}
	// Add reward from coinbase to miner
	ub.AddRewardTx(*bt, miner)
	// Build transactions ordered map
	txsOrderedMap := *orderedmap.NewOrderedMap()
	for f, t := range ub.Transactions {
		txsOrderedMap.Set(f, t)
	}

	return &MintedBlock{
		tree: bt,

		transactions: txsOrderedMap,
		previous:     ub.Previous,
		nonce:        0,
	}, nil
}

// Graft (mine and add) a block into a branch
func (bt *Blocktree) Graft(ub UnmintedBlock, miner string) Hash {
	// Get minted block
	b, err := bt.Mint(ub, miner)
	if err != nil {
		panic(err)
	}
	// Mine block
	b.mine(bt.network)
	hash := b.Hash()
	// Add to blocktree
	bt.blocks[hash] = *b
	return hash
}

// View tree
func (bt *Blocktree) View(branch Hash) {
	for i, b := 0, branch; b != NilHash; i, b = i+1, bt.blocks[b].previous {
		if bt.blocks[b].previous == NilHash {
			fmt.Printf("╳  #r\t%x\n", b)
		} else if i == 0 {
			fmt.Printf("┌─ #%d\t%x\n", i, b)
		} else {
			fmt.Printf("├─ #%d\t%x\n", i, b)
		}
	}
}

// Get list of unspent transactions
func (bt *Blocktree) findUnspentTxs(address string, branch Hash) []Transaction {
	var unspentTxs []Transaction
	spentTxIDs := map[Hash][]uint8{}

	for b := branch; b != NilHash; b = bt.blocks[b].previous {
		txs := bt.blocks[b].transactions
		for el := txs.Front(); el != nil; el = el.Next() {
			tx := el.Value.(Transaction)
			txHash := tx.Hash()
		Outputs:
			for i, out := range tx.Outputs {
				if spentTxIDs[txHash] != nil {
					for _, spentOut := range spentTxIDs[txHash] {
						if spentOut == uint8(i) {
							continue Outputs
						}
					}
				}
				if out.CanBeUnlocked(address) {
					unspentTxs = append(unspentTxs, tx)
				}
			}
			if !tx.IsCoinbase() {
				for _, in := range tx.Inputs {
					inTxID := in.ID
					spentTxIDs[txHash] = append(spentTxIDs[inTxID], in.Index)
				}
			}
		}
	}

	return unspentTxs
}

// Get unspent transactions output
func (tree *Blocktree) findUTxO(address string, branch Hash) []TxOutput {
	var UTXOs []TxOutput
	unspentTransactions := tree.findUnspentTxs(address, branch)
	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}
	return UTXOs
}

// Get spendable outputs for the provided wallet
func (tree *Blocktree) findSpendableOutputs(address string, amount uint64, branch Hash) (uint64, map[Hash][]uint8) {
	unspentOuts := map[Hash][]uint8{}
	unspentTxs := tree.findUnspentTxs(address, branch)
	accumulated := uint64(0)

Work:
	for _, tx := range unspentTxs {
		hash := tx.Hash()
		for i, out := range tx.Outputs {
			if out.CanBeUnlocked(address) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[hash] = append(unspentOuts[hash], uint8(i))

				if accumulated >= amount {
					break Work
				}
			}
		}
	}
	return accumulated, unspentOuts
}

// Find a block in the Blocktree
// Returns a pointer to the block (not the actual reference
// to the minted block) if found or nil otherwise
func (bt *Blocktree) FindBlock(hash Hash) *MintedBlock {
	b, ok := bt.blocks[hash]
	if !ok {
		return nil
	} else {
		return &b
	}
}
