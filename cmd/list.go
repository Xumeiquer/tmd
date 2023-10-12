/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/Xumeiquer/tmd/internal/logger"
	"github.com/Xumeiquer/tmd/internal/tclient"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zelenin/go-tdlib/client"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List information about conversations.",
	Long: `List command displays information about the conversations. A table will be 
shown as a result of this command listing users, chats, channels, or forums.`,
	Run: listCmdEx,
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.PersistentFlags().StringVarP(&filterDate, "filter-date", "", filterDate, "filter results by date")
	listCmd.PersistentFlags().StringVarP(&filterContent, "filter", "", filterContent, "filter results by content")
	listCmd.PersistentFlags().StringVarP(&filterContentRE, "filter-re", "", filterContentRE, "filter results by regular expresion")

	// Bind flags to configuration file
	viper.BindPFlag("tmd.cmds.filters.date", listCmd.PersistentFlags().Lookup("filter-date"))
	viper.BindPFlag("tmd.cmds.filters.regex", listCmd.PersistentFlags().Lookup("filter-re"))
	viper.BindPFlag("tmd.cmds.filters.text", listCmd.PersistentFlags().Lookup("filter"))
}

func listCmdEx(cmd *cobra.Command, args []string) {
	log := logger.GetLog(logLevel, logType, logTo)
	log = log.With("context", "listCmdEx")

	tg := tclient.NewTGClient()
	tg.Authenticate()
	defer tg.Stop()

	log.Debug("Telegram client authnticated")

	if chatLimit < 1 {
		log.Error("limit can not be lower than 1")
		os.Exit(1)
	}

	log.Debug("getting chats")
	chatList, err := tg.GetChats(chatLimit)
	if err != nil {
		log.Error(fmt.Sprintf("Err: %s\n", err.Error()))
		os.Exit(1)
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Id", "Channel", "Group", "Supergroup", "Chat", "Bot", "Secret", "Forum", "Username", "Title"})

	log.Debug("output processor initialized")
	log.Debug("processing messages")

	// TODO: Filter by the filter options
	for idx, chatId := range chatList.ChatIds {
		row := ChatRow{}

		chat, err := tg.GetChat(chatId)
		if err != nil {
			log.Error(fmt.Sprintf("retrieving chat %d information", chatId), "msg", err.Error())
		}

		row.Id = chat.Id
		switch chat.Type.ChatTypeType() {
		case client.TypeChatTypePrivate:
			row.Chat = true
			row.Title = chat.Title
			if len(chat.Title) > 50 {
				row.Title = chat.Title[:50] + "..."
			}

			privateRequest := chat.Type.(*client.ChatTypePrivate)
			userInfo, err := tg.GetUser(privateRequest.UserId)
			if err != nil {
				log.Error("unable to get user details", "msg", err.Error())
				continue
			}

			if userInfo.Usernames != nil && len(userInfo.Usernames.ActiveUsernames) > 0 {
				row.Username = userInfo.Usernames.ActiveUsernames[0]
			}
		case client.TypeChatTypeBasicGroup:
			row.Group = true
			row.Title = chat.Title
			if len(chat.Title) > 50 {
				row.Title = chat.Title[:50] + "..."
			}

			groupRequest := chat.Type.(*client.ChatTypeBasicGroup)
			groupInfo, err := tg.GetGroup(groupRequest.BasicGroupId)
			if err != nil {
				log.Error("unable to get group details", "msg", err.Error())
				continue
			}

			if groupInfo.UpgradedToSupergroupId != 0 {
				row.Supergroup = true
				superGroupInfo, err := tg.GetSupergroup(groupInfo.UpgradedToSupergroupId)
				if err != nil {
					log.Error("unable to get super group details", "msg", err.Error())
					continue
				}

				if superGroupInfo.Usernames != nil && len(superGroupInfo.Usernames.ActiveUsernames) > 0 {
					row.Username = superGroupInfo.Usernames.ActiveUsernames[0]
				}

				if superGroupInfo.IsChannel {
					row.Channel = true
				}

				if superGroupInfo.IsForum {
					row.Forum = true
				}

			}
		case client.TypeChatTypeSupergroup:
			row.Supergroup = true
			row.Title = chat.Title
			if len(chat.Title) > 50 {
				row.Title = chat.Title[:50] + "..."
			}

			supergroupRequest := chat.Type.(*client.ChatTypeSupergroup)
			superGroupInfo, err := tg.GetSupergroup(supergroupRequest.SupergroupId)
			if err != nil {
				log.Error("unable to get super group details", "msg", err.Error())
				continue
			}

			if superGroupInfo.Usernames != nil && len(superGroupInfo.Usernames.ActiveUsernames) > 0 {
				row.Username = superGroupInfo.Usernames.ActiveUsernames[0]
			}

			if superGroupInfo.IsChannel {
				row.Channel = true
			}

			if superGroupInfo.IsForum {
				row.Forum = true
			}
		case client.TypeChatTypeSecret:
			row.Secret = true
			row.Title = chat.Title
			if len(chat.Title) > 50 {
				row.Title = chat.Title[:50] + "..."
			}

			secretRequest := chat.Type.(*client.ChatTypeSecret)

			user, err := tg.GetUser(secretRequest.UserId)
			if err != nil {
				log.Error("unable to get user details", "msg", err.Error())
				continue
			}

			if len(user.Usernames.ActiveUsernames) > 0 {
				row.Username = user.Usernames.ActiveUsernames[0]
			}
		}
		// "#", "Id", "Channel", "Group", "Supergroup", "Chat", "Bot", "Secret", "Forum", "Username", "Title"
		t.AppendRow([]interface{}{idx + 1, row.Id, row.Channel, row.Group, row.Supergroup, row.Chat, row.Bot, row.Secret, row.Forum, row.Username, row.Title})
	}

	log.Debug("rendering output")
	t.Render()
}
