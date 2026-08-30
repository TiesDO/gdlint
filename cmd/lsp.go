package cmd

import (
	"context"
	"log"
	"os"

	"github.com/TiesDO/gdlint/lsp"
	"github.com/spf13/cobra"
)

var lspCmd = &cobra.Command{
	Use:   "lsp",
	Short: "starts the stdin/out lsp server",
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.New(os.Stderr, "[LSP] ", log.Ltime|log.Lshortfile)
		server := lsp.NewServer(os.Stdin, os.Stdout, logger)
		server.Run(context.Background())
	},
}

func init() {
	rootCmd.AddCommand(lspCmd)
}
