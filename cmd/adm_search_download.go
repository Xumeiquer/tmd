//go:build adm
// +build adm

package cmd

import (
	"os"
	"time"

	"github.com/Xumeiquer/tmd/internal/config"
	"github.com/Xumeiquer/tmd/internal/logger"
	"github.com/Xumeiquer/tmd/internal/tclient"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

// admSearchFilesCmd represents the list command
var admSearchFilesCmd = &cobra.Command{
	Use:   "files",
	Short: "Search for files in the the Telegram library cache",
	Run:   admSearchFilesCmdEx,
}

func init() {
	admSearchCmd.AddCommand(admSearchFilesCmd)
}

func admSearchFilesCmdEx(cmd *cobra.Command, args []string) {
	log := logger.GetLog(logLevel, logType, logTo)
	log = log.With("context", "admSearchFilesCmdEx")

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

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "File Id", "Message Id", "Added on", "Completed on", "In pause", "Downloaded from Chat ID"})

	log.Debug("output processor initialized")
	log.Debug("processing messages")

	for idx, file := range searchResults.Files {
		row := SearchResult{}

		row.FileId = file.FileId
		row.AddedOn = file.AddDate
		row.CompletedOn = file.CompleteDate
		row.IsPaused = file.IsPaused
		row.MessageId = file.Message.Id
		row.ChatId = file.Message.ChatId

		addedDate := time.Unix(int64(row.AddedOn), 0)
		completedDate := time.Unix(int64(row.CompletedOn), 0)

		t.AppendRow([]interface{}{idx + 1, row.FileId, row.MessageId, addedDate.Format(time.RFC822), completedDate.Format(time.RFC822), row.IsPaused, row.ChatId})
	}

	log.Debug("rendering output")
	t.Render()
}
