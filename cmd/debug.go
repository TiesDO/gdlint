package cmd

import "github.com/spf13/cobra"

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "debug utilities for rules",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(debugCmd)
}
