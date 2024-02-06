/*
Copyright © 2024 Jaume Martin <jaumartin@gmail.com>
*/
package cmd

import (
	"log/slog"
	"os"
	"time"

	"github.com/Xumeiquer/tmd/internal/tg"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

// listTopicCmd represents the list command
var listTopicCmd = &cobra.Command{
	Use:   "topics",
	Short: "Show iformation about topics inside of a channel.",
	Long: `Topics command displays information about the forum topics inside of a channel. A table will be 
shown as a result of this command listing the topics with its ID.`,
	Run: listTopicCmdEx,
}

func init() {
	listCmd.AddCommand(listTopicCmd)

	listTopicCmd.Flags().Int64VarP(&channelId, "chat", "", channelId, "channel (forum) id where to list topics from")
}

func listTopicCmdEx(cmd *cobra.Command, args []string) {
	if channelId == 0 {
		slog.Error("channel flag is required")
		os.Exit(1)
	}

	tgc := tg.NewTGClient()
	tgc.Authenticate()
	defer tgc.Stop()

	dateFilter, _ := time.Parse(dateParserLayout, filterDate)
	topicsList, err := tgc.GetForumTopics(channelId, dateFilter, filterContent)
	if err != nil {
		slog.Error("unable to read topics from channel", "msg", err.Error())
		os.Exit(1)
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Channel", "Message Thread Id", "General", "Closed", "Hidden", "Name", "Creation Date"})

	for idx, topic := range topicsList {
		row := TopicRow{}
		row.Id = topic.Info.MessageThreadId
		row.Name = topic.Info.Name
		row.CreationDate = topic.Info.CreationDate
		row.General = topic.Info.IsGeneral
		row.Closed = topic.Info.IsClosed
		row.Hidden = topic.Info.IsHidden

		// "#", "Channel", "Id", "General", "Closed", "Hidden", "Name"
		date := time.Unix(int64(row.CreationDate), 0)
		t.AppendRow([]interface{}{idx + 1, channelId, row.Id, row.General, row.Closed, row.Hidden, row.Name, date.Format(time.RFC822)})
	}

	t.Render()
}
