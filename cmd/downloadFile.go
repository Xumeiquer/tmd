/*
Copyright © 2024 Jaume Martin <jaumartin@gmail.com>
*/
package cmd

import (
	"log/slog"
	"os"

	"github.com/Xumeiquer/tmd/internal/tg"
	"github.com/spf13/cobra"
)

// downloadFileCmd download media found on a message
var downloadFileCmd = &cobra.Command{
	Use:   "file",
	Short: "Download document",
	Long:  `Download media document from the Telegram cloud`,
	Run:   downloadFileCmdEx,
}

func init() {
	downloadCmd.AddCommand(downloadFileCmd)

	downloadFileCmd.Flags().Int64VarP(&chatId, "chat", "", chatId, "chat ID (chat/group/channel/...)")
	downloadFileCmd.Flags().StringVarP(&fileId, "file", "", fileId, "file ID to download")
}

func downloadFileCmdEx(cmd *cobra.Command, args []string) {
	tgc := tg.NewTGClient(tg.WithStorePath(storePath), tg.WithPrintStatus(true))
	tgc.Authenticate()
	defer tgc.Stop()

	listener := tgc.GetListener()
	dm := tg.NewDM(listener)
	tgc.SetDownloadManager(dm)

	_, err := tgc.DownloadFile(chatId, fileId)
	if err != nil {
		slog.Error("unable to download media", "msg", err.Error())
		os.Exit(1)
	}

	if dm.IsDownloadCompleted() {
		slog.Info("Download already completed")
		return
	}

	dm.Wait()
}
