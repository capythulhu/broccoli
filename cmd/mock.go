package cmd

import (
	"log"
	"math"
	"math/rand"
	"os"

	"github.com/spf13/cobra"
	"github.com/thzoid/broccoli/blocktree"
)

var dlog *log.Logger = log.New(os.Stdout, "", 0)
var wlog *log.Logger = log.New(os.Stdout, "warning: ", 0)

// mockCmd represents the test command
var mockCmd = &cobra.Command{
	Use:   "mock",
	Short: "run a mock Blocktree with random data",
	Long: `create and run a sample Blocktree with randomly
	generated data for each block. can be used for quickly
	setting up a mock environment.`,
	Run: func(cmd *cobra.Command, args []string) {
		difficulty, _ := cmd.PersistentFlags().GetUint8("difficulty")
		if difficulty > 22 {
			wlog.Println("difficulty is higher than 22. blocks might take several minutes to be mined.")
		}
		branches, _ := cmd.PersistentFlags().GetUint16("branches")
		blocks, _ := cmd.PersistentFlags().GetUint16("blocks")
		if blocks == 0 {
			wlog.Println("branches cannot have less than 1 block. defaulting to 1.")
			blocks = 1
		}

		tree, root := blocktree.NewTree(&blocktree.Network{Difficulty: difficulty}, "root-address")
		dlog.Printf("root: %x", root)

		for i := uint16(0); i < branches; i++ {
			dlog.Printf("branch %d:\n", i)
			blocksOnBranch := uint16(math.Ceil(float64(rand.Float32() * float32(blocks))))
			previousBlock := root
			for j := uint16(0); j < blocksOnBranch; j++ {
				previousBlock = tree.NewBlock([]*blocktree.Transaction{}, previousBlock)
			}
			tree.View(previousBlock)
		}

	},
}

func init() {
	rootCmd.AddCommand(mockCmd)
	mockCmd.PersistentFlags().Uint16P("branches", "b", 3, "number of branches to be generated")
	mockCmd.PersistentFlags().Uint16("blocks", 5, "maximum number of blocks to be generated on each branch")
	mockCmd.PersistentFlags().Uint8P("difficulty", "d", 12, "mining difficulty (amount of zeroed bits on the left side of resulting block hashes)")

}
