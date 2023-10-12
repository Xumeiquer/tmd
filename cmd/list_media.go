package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/Xumeiquer/tmd/internal/logger"
	"github.com/Xumeiquer/tmd/internal/tclient"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zelenin/go-tdlib/client"
)

// listMediaCmd represents the list command
var listMediaCmd = &cobra.Command{
	Use:   "media",
	Short: "List media displais the media file names",
	Long:  `List media command displays information about the media files store on a message, chat, channel, topic, etc.`,
	Run:   listMediaCmdEx,
}

func init() {
	listCmd.AddCommand(listMediaCmd)

	listMediaCmd.Flags().Int64VarP(&chatId, "chat", "", chatId, "chat/channel ID where to list media from")
	listMediaCmd.Flags().Int64VarP(&topicId, "topic", "", topicId, "topic ID where to list media from")
	listMediaCmd.Flags().StringVarP(&messageUrl, "message-url", "", messageUrl, "url to a particular message")

	// Bind flags to configuration file
	viper.BindPFlag("tmd.cmds.list-media.chat", listMediaCmd.PersistentFlags().Lookup("chat"))
	viper.BindPFlag("tmd.cmds.list-media.topic", listMediaCmd.PersistentFlags().Lookup("topic"))
	viper.BindPFlag("tmd.cmds.list-media.message", listMediaCmd.PersistentFlags().Lookup("message-url"))
}

func listMediaCmdEx(cmd *cobra.Command, args []string) {
	log := logger.GetLog(logLevel, logType, logTo)
	log = log.With("context", "listMediaCmdEx")

	if filterDate == initialDate {
		log.Warn("defining a date is not mandatory, but it can reduce the amount of messages to retrieve.")
	}

	tg := tclient.NewTGClient()
	tg.Authenticate()
	defer tg.Stop()

	log.Debug("Telegram client authnticated")

	var messages []*client.Message
	var err error

	dateFilter, _ := time.Parse(dateParserLayout, filterDate)

	log.Debug("getting chat history")
	if messageUrl != "" {
		log.Debug(fmt.Sprintf("using message url: %s", messageUrl))
		linkInfo, err := tg.GetMessageLinkInfo(messageUrl)
		if err != nil {
			log.Error("unable to read the message", "msg", err.Error())
			os.Exit(1)
		}
		messages, err = tg.GetForumHistory(linkInfo.ChatId, linkInfo.MessageThreadId, dateFilter, filterContent, filterContentRE)
		if err != nil {
			log.Error("unable to get message media", "msg", err.Error())
			os.Exit(1)
		}
	} else if topicId != 0 {
		if chatId == 0 {
			log.Error(fmt.Sprintf("chat ID is mandatory"))
			os.Exit(1)
		}
		log.Debug(fmt.Sprintf("using chat ID: %d topid ID: %d", chatId, topicId))
		messages, err = tg.GetForumHistory(chatId, topicId, dateFilter, filterContent, filterContentRE)
		if err != nil {
			log.Error("unable to get topic", "msg", err.Error())
			os.Exit(1)
		}
	} else {
		if chatId == 0 {
			log.Error(fmt.Sprintf("chat ID is mandatory"))
			os.Exit(1)
		}
		log.Debug(fmt.Sprintf("using chat ID: %d", chatId))
		messages, err = tg.GetChatHistory(chatId, dateFilter, filterContent, filterContentRE)
		if err != nil {
			log.Error("unable to get messages", "msg", err.Error())
			os.Exit(1)
		}
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Chat Id", "Message Id", "File Id", "Name", "Mime", "Size", "Date"})

	log.Debug("output processor initialized")
	log.Debug("processing messages")

	for idx, message := range messages {
		row := DocumentRow{}

		if message == nil {
			log.Warn("nil message found")
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

	log.Debug("rendering output")
	t.Render()
}

func getMessageContent(msg *client.Message, length int) string {
	switch msg.GetType() {
	case client.TypeMessage:
		content, ok := msg.Content.(*client.MessageText)
		if !ok {
			return ""
		}
		return content.Text.Text
	case client.TypeDocument:
		content, ok := msg.Content.(*client.MessageDocument)
		if !ok {
			return ""
		}
		if len(content.Caption.Text) > length {
			return fmt.Sprintf("%s (%s) [%s]", content.Document.FileName, content.Document.Type, content.Caption.Text[:length])
		}
		return fmt.Sprintf("%s (%s) [%s]", content.Document.FileName, content.Document.Type, content.Caption.Text)
	}
	return ""
}
