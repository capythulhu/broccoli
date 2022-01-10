package main

import (
	"fmt"

	"github.com/thzoid/broccoli/blocktree"
	"github.com/thzoid/broccoli/hash"
	"github.com/thzoid/broccoli/wallet"
)

func BlockView(bt *blocktree.Blocktree, h hash.Hash) {
	block := bt.Block(h)
	if block.Previous() == hash.NilHash {
		fmt.Printf("block %s (root)\n", h.String())
	} else {
		fmt.Printf("block %s\n", h.String())
	}
	for sender, tx := range block.Transactions() {
		var wallet string
		if sender.IsCoinbase() {
			wallet = "coinbase"
		} else {
			wallet = sender.String()
		}
		fmt.Printf("├─tx %s\t\t ◀ %s\n", tx.Hash().String(), wallet)
		txConnector := '├'
		for j, in := range tx.Inputs {
			if j == len(tx.Inputs)-1 && len(tx.Outputs) == 0 {
				txConnector = '└'
			}
			fmt.Printf("│ %c─in  %s[%d]\n", txConnector, in.ID.String(), in.Index)
		}
		noChange := tx.Outputs[1].Value == 0
		for j, out := range tx.Outputs {
			if j == 0 && noChange || j == 1 {
				txConnector = '└'
			}
			if out.Value > 0 {
				fmt.Printf("│ %c─out %d\t\t ▶ %s\n", txConnector, out.Value, out.Address.String())
			}
		}
		fmt.Printf("│\n")
	}
}

func main() {
	fmt.Println(
		`
		BROCCOLI 1.0.0

                      ████              
            ██████  ██░░░░██  ████      
          ██░░░░░░██░░░░░░░░██░░░░██    
        ██░░░░░░░░░░░░░░░░░░░░░░░░▒▒██  
    ████░░░░░░░░░░░░░░░░░░░░░░░░▒▒▒▒▒▒██
  ██▒▒░░░░░░░░░░░░░░░░░░░░▒▒░░▒▒▒▒▒▒▒▒██
  ██▒▒▒▒░░░░░░░░░░░░░░░░░░▒▒▒▒▒▒▒▒▒▒██  
  ██▒▒▒▒▒▒░░░░▒▒░░▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒██    
    ██▒▒▒▒▒▒░░▒▒▒▒▒▒████▒▒██▒▒▒▒██      
    ██▒▒▒▒▒▒▒▒▒▒▒▒██░░░░██  ██░░██      
      ██▒▒▒▒██████  ██░░░░██░░██        
        ████  ██░░██░░░░░░░░░░██        
                ██░░░░░░░░░░░░██        
                  ██░░░░░░░░░░██        
                  ██░░░░░░░░░░██        
                  ██░░░░░░░░░░██        
                  ██░░░░░░░░░░██        
                    ██░░░░░░██          
                      ██████            

					  `)

	wallets := make([]wallet.Wallet, 10)
	for i := range wallets {
		wallets[i] = wallet.NewWallet()
	}

	alice := wallets[0].Address()
	bob := wallets[1].Address()
	carol := wallets[2].Address()

	fmt.Println("alice\t" + alice.String())
	fmt.Println("bob\t" + bob.String())
	fmt.Println("carol\t" + carol.String())
	fmt.Print("\n")

	tree, root := blocktree.NewTree(blocktree.Network{Difficulty: 1, Reward: 100}, alice)
	BlockView(&tree, root)

	b1 := blocktree.NewBlock(root)
	b1.AddTx(tree, alice, blocktree.TxOutput{Address: bob, Value: 5})
	b1Hash := tree.Graft(b1, carol)
	BlockView(&tree, b1Hash)

	// fmt.Println("\nTree view")
	// tree.View(b3Hash)
}
