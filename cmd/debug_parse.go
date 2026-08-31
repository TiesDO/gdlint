package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/TiesDO/gdlint/rules"
	"github.com/spf13/cobra"
)

var debugParseCmd = &cobra.Command{
	Use:   "parse [filepath]",
	Short: "output the tree sitter nodes for a file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pathStr := args[0]

		path, err := filepath.Abs(pathStr)

		if err != nil {
			fmt.Printf("failed to expand provided filepath: %v\n", err)
			os.Exit(1)
		}

		content, err := os.ReadFile(path)

		if err != nil {
			fmt.Printf("failed to read file: %v", err)
			os.Exit(1)
		}

		runner := rules.NewRuleRunner(&rules.DefaultRuleRegistry, content)

		runner.UpdateTree(context.Background())
		fmt.Printf("%s\n\n", runner.SprintTree())
	},
}

func init() {
	debugCmd.AddCommand(debugParseCmd)
}
