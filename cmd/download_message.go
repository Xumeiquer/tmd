package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/Xumeiquer/tmd/internal/logger"
	progress_bar "github.com/Xumeiquer/tmd/internal/progress-bar"
	"github.com/Xumeiquer/tmd/internal/tclient"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zelenin/go-tdlib/client"
)

// downloadMessageCmd download media found on a message
var downloadMessageCmd = &cobra.Command{
	Use:   "message",
	Short: "Download media from a Telegram message.",
	Long:  `Download message reads a Telegram message and downloads media files linked in the message.`,
	Run:   downloadMessageEx,
}

func init() {
	downloadCmd.AddCommand(downloadMessageCmd)

	downloadMessageCmd.Flags().StringVarP(&messageUrl, "message-url", "", messageUrl, "message url where to download media from")

	// Bind flags to configuration file
	viper.BindPFlag("tmd.cmds.download-message.message", downloadMessageCmd.PersistentFlags().Lookup("message-url"))
}

func downloadMessageEx(cmd *cobra.Command, args []string) {
	log := logger.GetLog(logLevel, logType, logTo)
	log = log.With("context", "downloadMessageEx")

	if messageUrl == "" {
		log.Error("message URL cant't be empty")
		os.Exit(1)
	}

	tg := tclient.NewTGClient()
	tg.Authenticate()
	defer tg.Stop()

	linkInfo, err := tg.GetMessageLinkInfo(messageUrl)
	if err != nil {
		log.Error(fmt.Sprintf("unable to get the message. Err: %s", err.Error()))
	}

	dateFilter, _ := time.Parse(dateParserLayout, filterDate)
	messages, err := tg.GetForumHistory(linkInfo.ChatId, linkInfo.MessageThreadId, dateFilter, filterContent, filterContentRE)
	if err != nil {
		log.Error("unable to read message", "msg", err.Error())
		os.Exit(1)
	}

	listener := tg.GetListener()
	defer listener.Close()

	log.Debug("got Telegram listener")
	log.Debug("listener ready to spwan")

	sync := tg.GetSynchChannel()

	var prb *progress_bar.ProgressBar = nil
	if showProgress {
		prb = progress_bar.New()
	}

	go tg.Poll(sync, listener, prb)

	for _, message := range messages {
		log.Debug(fmt.Sprintf("message type: %s", message.GetType()))

		switch message.Content.MessageContentType() {
		case client.TypeMessageDocument:
			content := message.Content.(*client.MessageDocument)
			log.Debug(fmt.Sprintf("download request for file %d in chat %d", content.Document.Document.Id, linkInfo.ChatId))
			file, err := tg.DownloadFile(linkInfo.ChatId, content.Document.Document.Remote.Id)
			if err != nil {
				log.Error("unable to download media", "msg", err.Error())
				continue
			}

			if prb != nil {
				size := file.Size
				if file.Size == 0 {
					size = file.ExpectedSize
					if file.ExpectedSize == 0 {
						size = file.Remote.UploadedSize
					}
				}

				log.Debug(fmt.Sprintf("adding file %d to the tracker", file.Id))
				prb.AddTracker(file.Id, size)
			}
		}
	}

	log.Debug("waiting for download to be completed")
	tg.Wait(sync)

	tg.MoveDownloaded(storePath)
}
