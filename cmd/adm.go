//go:build adm
// +build adm

package cmd

import (
	"github.com/spf13/cobra"
)

// admCmd represents the list command
var admCmd = &cobra.Command{
	Use:   "adm",
	Short: "Run administrative tasks",
	Long: `Run a small set of administrative tasks such as remove downloaded files from the cache,
search for downloaded media, etc.`,
}

func init() {
	rootCmd.AddCommand(admCmd)
}
