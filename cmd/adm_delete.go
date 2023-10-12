//go:build adm
// +build adm

package cmd

import (
	"github.com/spf13/cobra"
)

// admDeleteCmd represents the list command
var admDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete media files from the Telegram cache library",
}

func init() {
	admCmd.AddCommand(admDeleteCmd)
}
