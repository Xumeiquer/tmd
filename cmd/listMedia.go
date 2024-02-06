/*
Copyright © 2024 Jaume Martin <jaumartin@gmail.com>
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Xumeiquer/tmd/internal/tg"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/zelenin/go-tdlib/client"
)

// listMediaCmd represents the list command
var listMediaCmd = &cobra.Command{
	Use:   "media",
	Short: "Displais the media files",
	Long:  `Media command displays information about the media files stored on a messages, chats, channels, topics, etc.`,
	Run:   listMediaCmdEx,
}

func init() {
	listCmd.AddCommand(listMediaCmd)

	listMediaCmd.Flags().Int64VarP(&chatId, "chat", "", chatId, "chat/channel ID where to list media from")
	listMediaCmd.Flags().Int64VarP(&topicId, "topic", "", topicId, "topic ID where to list media from")
	listMediaCmd.Flags().StringVarP(&messageUrl, "message-url", "", messageUrl, "url to a particular message")
}

func listMediaCmdEx(cmd *cobra.Command, args []string) {
	if filterDate == initialDate {
		slog.Warn("defining a date is not mandatory, but it can reduce the amount of messages to retrieve.")
	}

	tgc := tg.NewTGClient()
	tgc.Authenticate()
	defer tgc.Stop()

	dateFilter, _ := time.Parse(dateParserLayout, filterDate)

	var messages []*client.Message
	var err error

	slog.Debug("getting chat history")
	if messageUrl != "" {
		slog.Debug(fmt.Sprintf("using message url: %s", messageUrl))
		linkInfo, err := tgc.GetMessageLinkInfo(messageUrl)
		if err != nil {
			slog.Error("unable to read the message", "msg", err.Error())
			os.Exit(1)
		}
		messages, err = tgc.GetForumHistory(linkInfo.ChatId, linkInfo.MessageThreadId, dateFilter, filterContent)
		if err != nil {
			slog.Error("unable to get message media", "msg", err.Error())
			os.Exit(1)
		}
	} else if topicId != 0 {
		if chatId == 0 {
			slog.Error(fmt.Sprintf("chat ID is mandatory"))
			os.Exit(1)
		}
		slog.Debug(fmt.Sprintf("using chat ID: %d topid ID: %d", chatId, topicId))
		messages, err = tgc.GetForumHistory(chatId, topicId, dateFilter, filterContent)
		if err != nil {
			slog.Error("unable to get topic", "msg", err.Error())
			os.Exit(1)
		}
	} else {
		if chatId == 0 {
			slog.Error(fmt.Sprintf("chat ID is mandatory"))
			os.Exit(1)
		}
		slog.Debug(fmt.Sprintf("using chat ID: %d", chatId))
		messages, err = tgc.GetChatHistory(chatId, dateFilter, filterContent)
		if err != nil {
			slog.Error("unable to get messages", "msg", err.Error())
			os.Exit(1)
		}
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Chat Id", "Message Id", "File Id", "Name", "Mime", "Size", "Date"})

	for idx, message := range messages {
		row := DocumentRow{}

		if message == nil {
			continue
		}

		if mediaInfo, ok := message.Content.(*client.MessageDocument); ok {
			row.Id = mediaInfo.Document.Document.Remote.Id
			row.ChatId = message.ChatId
			row.MessageId = message.Id
			// row.Message = getMessageContent(message, 50)
			row.Name = mediaInfo.Document.FileName
			row.Mime = mediaInfo.Document.MimeType
			row.Date = message.Date

			if mediaInfo.Document.Document.Size != 0 {
				row.Size = mediaInfo.Document.Document.Size
			} else if mediaInfo.Document.Document.ExpectedSize != 0 {
				row.Size = mediaInfo.Document.Document.ExpectedSize
			} else {
				row.Size = 0
			}
		}
		// "#", "Chat Id", "Message Id", "File Id", "Name", "Mime", "Size", "Date"
		date := time.Unix(int64(row.Date), 0)
		t.AppendRow([]interface{}{idx + 1, row.ChatId, row.MessageId, row.Id, row.Name, row.Mime, row.Size, date.Format(time.RFC822)})

	}

	t.Render()
}
