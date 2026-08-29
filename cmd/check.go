package cmd

import (
	"context"
	"fmt"
	"os"

	"path/filepath"

	"github.com/TiesDO/gdlint/rules"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check [filepath]",
	Short: "Runs the checks against a file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pathStr := args[0]

		path, err := filepath.Abs(pathStr)

		if err != nil {
			fmt.Printf("failed to expand provided filepath: %v", err)
			os.Exit(1)
		}

		content, err := os.ReadFile(path)

		if err != nil {
			fmt.Printf("failed to read file: %v", err)
			os.Exit(1)
		}

		runner := rules.NewRuleRunner(&rules.DefaultRuleRegistry, content)

		// TODO: add arguments to pass rules you want
		warnings, err := runner.RunRules([]string{"untyped_function_argument", "untyped_variable_statement", "untyped_function_return"}, context.Background())

		if err != nil {
			fmt.Printf("error while executing rules: %v", err)
			os.Exit(1)
		}

		if len(warnings) == 0 {
			fmt.Println("no warnings found!")
			os.Exit(0)
		}

		fmt.Printf("found %d warnings\n", len(warnings))
		for _, warning := range warnings {
			fmt.Printf("line %d: %s (@%s)\n", warning.LineNumber, warning.Message, warning.Offense)
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
