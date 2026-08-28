package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var dumpCmd = &cobra.Command{
	Use:   "dump [filepath]",
	Short: "Dumps interpreted treesitter syntax",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filepath := args[0]
		fmt.Printf("dump command on file %s\n", filepath)
	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)
}
