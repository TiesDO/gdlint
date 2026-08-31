package cmd

import (
	"context"
	"fmt"
	"os"
	"slices"

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
			fmt.Printf("failed to expand provided filepath: %v\n", err)
			os.Exit(1)
		}

		content, err := os.ReadFile(path)

		if err != nil {
			fmt.Printf("failed to read file: %v", err)
			os.Exit(1)
		}

		runner := rules.NewRuleRunner(&rules.DefaultRuleRegistry, content)

		included_rules, err := cmd.Flags().GetStringSlice("include")

		if err != nil {
			fmt.Printf("failed to extract included rules: %v\n", err)
			os.Exit(1)
		}

		if len(included_rules) == 0 {
			included_rules = rules.DefaultRuleRegistry.RuleNames()
		}

		excluded_rules, err := cmd.Flags().GetStringSlice("exclude")

		if err != nil {
			fmt.Printf("failed to extract excluded rules: %v\n", err)
			os.Exit(1)
		}

		target_rules := make([]string, 0)

		if len(excluded_rules) > 0 {
			for _, included_rule := range included_rules {
				is_included := !slices.Contains(excluded_rules, included_rule)

				if is_included {
					target_rules = append(target_rules, included_rule)
				}
			}
		} else {
			target_rules = included_rules
		}

		warnings, err := runner.RunRules(target_rules, context.Background())

		if err != nil {
			fmt.Printf("error while executing rules: %v\n", err)
			os.Exit(1)
		}

		if len(warnings) == 0 {
			fmt.Println("no warnings found!")
			os.Exit(0)
		}

		fmt.Printf("found %d warnings\n", len(warnings))
		for _, warning := range warnings {
			fmt.Printf("%s\n", warning.FullMessage())
		}
	},
}

func init() {
	checkCmd.Flags().StringSliceP("include", "i", []string{}, "rules to include")
	checkCmd.Flags().StringSliceP("exclude", "e", []string{}, "rules to exclude (on conflict overrides include)")
	rootCmd.AddCommand(checkCmd)
}
