//go:build adm
// +build adm

package cmd

import (
	"github.com/Xumeiquer/tmd/internal/config"
	"github.com/Xumeiquer/tmd/internal/logger"
	"github.com/Xumeiquer/tmd/internal/tclient"
	"github.com/spf13/cobra"
)

// admMoveCmd represents the list command
var admMoveCmd = &cobra.Command{
	Use:   "move",
	Short: "Move command will move files",
	Long:  `Move command will move files from the download directory to the final store path`,
	Run:   admMoveCmdEx,
}

func init() {
	admCmd.AddCommand(admMoveCmd)
	admMoveCmd.Flags().StringVarP(&storePath, "store", "", storePath, "final store path to move downloaded files")
}

func admMoveCmdEx(cmd *cobra.Command, args []string) {
	log := logger.GetLog(logLevel, logType, logTo)
	log = log.With("context", "admMoveCmdEx")

	cfg := config.New()

	tg := tclient.NewTGClient(cfg)
	tg.Authenticate()
	defer tg.Stop()

	tg.MoveDownloaded(storePath)
}
