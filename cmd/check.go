package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check [filepath]",
	Short: "Runs the checks against a file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filepath := args[0]

		fmt.Printf("check command on %s\n", filepath)
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
