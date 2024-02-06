/*
Copyright © 2024 Jaume Martin <jaumartin@gmail.com>
*/
package cmd

import (
	"log/slog"
	"os"
	"time"

	"github.com/Xumeiquer/tmd/internal/tg"
	"github.com/spf13/cobra"
	"github.com/zelenin/go-tdlib/client"
)

// downloadMessageCmd download media found on a message
var downloadMessageCmd = &cobra.Command{
	Use:   "message",
	Short: "Download all files from a message",
	Long:  `Download all media document stored in a message from the Telegram cloud`,
	Run:   downloadMessageCmdEx,
}

func init() {
	downloadCmd.AddCommand(downloadMessageCmd)

	downloadMessageCmd.Flags().Int64VarP(&chatId, "chat", "", chatId, "chat ID (chat/group/channel/...)")
	downloadMessageCmd.Flags().StringVarP(&messageUrl, "message-url", "", messageUrl, "topic ID to download media from")
}

func downloadMessageCmdEx(cmd *cobra.Command, args []string) {
	tgc := tg.NewTGClient(tg.WithStorePath(storePath), tg.WithPrintStatus(true))
	tgc.Authenticate()
	defer tgc.Stop()

	listener := tgc.GetListener()
	dm := tg.NewDM(listener)
	tgc.SetDownloadManager(dm)

	dateFilter, _ := time.Parse(dateParserLayout, filterDate)

	linkInfo, err := tgc.GetMessageLinkInfo(messageUrl)
	if err != nil {
		slog.Error("unable to read the message", "msg", err.Error())
		os.Exit(1)
	}

	messages, err := tgc.GetForumHistory(linkInfo.ChatId, linkInfo.MessageThreadId, dateFilter, filterContent)
	if err != nil {
		slog.Error("unable to get topic", "msg", err.Error())
		os.Exit(1)
	}

	for _, message := range messages {
		if mediaInfo, ok := message.Content.(*client.MessageDocument); ok {
			fileId := mediaInfo.Document.Document.Remote.Id

			_, err := tgc.DownloadFile(chatId, fileId)
			if err != nil {
				slog.Error("unable to download media", "msg", err.Error())
				os.Exit(1)
			}
		}
	}

	if dm.IsDownloadCompleted() {
		slog.Info("Download already completed")
		return
	}

	dm.Wait()
}
