package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/TiesDO/gdlint/rules"
	"github.com/spf13/cobra"
)

var debugQueryCmd = &cobra.Command{
	Use:   "query [filepath]",
	Short: "run a tree-sitter query against a file",
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

		pattern, err := cmd.Flags().GetString("pattern")

		if err != nil {
			fmt.Printf("failed to extract pattern flag: %v", err)
			os.Exit(1)
		}

		parser := rules.NewParser()
		err = parser.Parse(context.Background(), content)

		if err != nil {
			fmt.Printf("failed to parse source into tree: %v", err)
			os.Exit(1)
		}

		output, err := parser.SQuery(pattern, content)

		if err != nil {
			fmt.Printf("failed to run query: %v", err)
			os.Exit(1)
		}

		fmt.Printf("%s\n", output)
	},
}

func init() {
	debugQueryCmd.Flags().StringP("pattern", "p", "", "the pattern the query the file with")
	debugCmd.AddCommand(debugQueryCmd)
}
