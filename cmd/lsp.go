package cmd

import (
	"context"
	"fmt"
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
		err := server.Run(context.Background())

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(lspCmd)
}
