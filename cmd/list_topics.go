package cmd

import (
	"os"
	"time"

	"github.com/Xumeiquer/tmd/internal/logger"
	"github.com/Xumeiquer/tmd/internal/tclient"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listTopicCmd represents the list command
var listTopicCmd = &cobra.Command{
	Use:   "topics",
	Short: "List topic (forum) information.",
	Long: `List topic command displays information about the forum topics inside a channel. A table will be 
shown as a result of this command listing users, chats, channels, or forums.`,
	Run: listTopicCmdEx,
}

func init() {
	listCmd.AddCommand(listTopicCmd)

	listTopicCmd.Flags().Int64VarP(&channelId, "channel", "", channelId, "channel (forum) id where to list topics from")

	// Bind flags to configuration file
	viper.BindPFlag("tmd.cmds.list-topics.channel", listTopicCmd.PersistentFlags().Lookup("channel"))
}

func listTopicCmdEx(cmd *cobra.Command, args []string) {
	log := logger.GetLog(logLevel, logType, logTo)
	log = log.With("context", "listTopicCmdEx")

	if channelId == 0 {
		log.Error("channel flag is required")
		os.Exit(1)
	}

	tg := tclient.NewTGClient()
	tg.Authenticate()
	defer tg.Stop()

	dateFilter, _ := time.Parse(dateParserLayout, filterDate)
	topicsList, err := tg.GetForumTopics(channelId, dateFilter, filterContent, filterContentRE)
	if err != nil {
		log.Error("unable to read topics from channel", "msg", err.Error())
		os.Exit(1)
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Channel", "Message Thread Id", "General", "Closed", "Hidden", "Name", "Creation Date"})

	log.Debug("output processor initialized")
	log.Debug("processing topics")

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

	log.Debug("rendering output")
	t.Render()
}
