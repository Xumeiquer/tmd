package cmd

import (
	"fmt"
	"os"

	"github.com/Xumeiquer/tmd/internal/logger"
	progress_bar "github.com/Xumeiquer/tmd/internal/progress-bar"
	"github.com/Xumeiquer/tmd/internal/tclient"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	// Bind flags to configuration file
	viper.BindPFlag("tmd.cmds.download-file.chat", downloadFileCmd.PersistentFlags().Lookup("chat"))
	viper.BindPFlag("tmd.cmds.download-file.file", downloadFileCmd.PersistentFlags().Lookup("file"))
}

func downloadFileCmdEx(cmd *cobra.Command, args []string) {
	log := logger.GetLog(logLevel, logType, logTo)
	log = log.With("context", "downloadFileCmdEx")

	tg := tclient.NewTGClient()
	tg.Authenticate()
	defer tg.Stop()

	log.Debug("Telegram client authnticated")

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

	file, err := tg.DownloadFile(chatId, fileId)
	if err != nil {
		log.Error("unable to download media", "msg", err.Error())
		os.Exit(1)
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

	log.Debug("waiting for download to be completed")
	tg.Wait(sync)

	tg.MoveDownloaded(storePath)
}

/*
 ch := make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		tdlibClient.Stop()
		os.Exit(1)
	}()
*/
