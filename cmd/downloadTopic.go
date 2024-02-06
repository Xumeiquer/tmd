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

// downloadTopicCmd download media found on a message
var downloadTopicCmd = &cobra.Command{
	Use:   "topic",
	Short: "Download all files in a topic or forum",
	Long:  `Download all media document stored in a topic or channel forum from the Telegram cloud`,
	Run:   downloadTopicCmdEx,
}

func init() {
	downloadCmd.AddCommand(downloadTopicCmd)

	downloadTopicCmd.Flags().Int64VarP(&chatId, "chat", "", chatId, "chat ID (chat/group/channel/...)")
	downloadTopicCmd.Flags().Int64VarP(&topicId, "topic", "", topicId, "topic ID to download media from")
}

func downloadTopicCmdEx(cmd *cobra.Command, args []string) {
	tgc := tg.NewTGClient(tg.WithStorePath(storePath), tg.WithPrintStatus(true))
	tgc.Authenticate()
	defer tgc.Stop()

	listener := tgc.GetListener()
	dm := tg.NewDM(listener)
	tgc.SetDownloadManager(dm)

	dateFilter, _ := time.Parse(dateParserLayout, filterDate)

	messages, err := tgc.GetForumHistory(chatId, topicId, dateFilter, filterContent)
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
