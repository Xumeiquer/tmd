//go:build adm
// +build adm

package cmd

import (
	"fmt"
	"os"

	"github.com/Xumeiquer/tmd/internal/config"
	"github.com/Xumeiquer/tmd/internal/logger"
	"github.com/Xumeiquer/tmd/internal/tclient"
	"github.com/spf13/cobra"
)

// admDeleteFilesCmd represents the list command
var admDeleteFilesCmd = &cobra.Command{
	Use:   "files",
	Short: "Delete media files from the Telegram cache library",
	Run:   admDeleteFilesCmdEx,
}

func init() {
	admDeleteCmd.AddCommand(admDeleteFilesCmd)
}

func admDeleteFilesCmdEx(cmd *cobra.Command, args []string) {
	log := logger.GetLog(logLevel, logType, logTo)
	log = log.With("context", "admDeleteFilesCmdEx")

	cfg := config.New()

	tg := tclient.NewTGClient(cfg)
	tg.Authenticate()
	defer tg.Stop()

	log.Debug("Telegram client authnticated")

	searchResults, err := tg.SearchFileDownloads()
	if err != nil {
		log.Error("unable to search for downloaded media", "msg", err.Error())
		os.Exit(1)
	}

	if len(searchResults.Files) == 0 {
		log.Info("There are no files to delete.")
	} else {
		for _, file := range searchResults.Files {
			ok, err := tg.DeleteFile(file.FileId)
			if err != nil {
				log.Error(fmt.Sprintf("unable to delete file ID %d", file.FileId), "msg", err.Error())
			}
			if ok != nil {
				log.Debug(fmt.Sprintf("Type: %s || Class: %s", ok.GetType(), ok.GetClass()))
			}
		}
	}
}
