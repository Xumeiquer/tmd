//go:build adm
// +build adm

package cmd

import (
	"github.com/spf13/cobra"
)

// admSearchCmd represents the list command
var admSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search on the Telegram library cache",
	Long:  `Search for different items stored by Telegram in the library`,
}

func init() {
	admCmd.AddCommand(admSearchCmd)
}
