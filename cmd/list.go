/*
Copyright © 2024 Jaume Martin <jaumartin@gmail.com>
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/Xumeiquer/tmd/internal/tg"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
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

	listCmd.PersistentFlags().StringVarP(&filterDate, "date", "", filterDate, "filter results by date")
	listCmd.PersistentFlags().StringVarP(&filterContent, "filter", "", filterContent, "filter results by content")
	listCmd.PersistentFlags().StringVarP(&filterType, "type", "", filterType, "filter results by chat type (group, supergroup, chat, secret)")
}

func listCmdEx(cmd *cobra.Command, args []string) {
	tgc := tg.NewTGClient()
	tgc.Authenticate()
	defer tgc.Stop()

	if chatLimit < 1 {
		slog.Error("limit can not be lower than 1")
		os.Exit(1)
	}

	types := []string{}
	justSuperGroups := false
	terms := strings.Split(filterType, ",")
	for _, t := range terms {
		typ := strings.Trim(t, " ")
		switch typ {
		case "chat":
			types = append(types, client.TypeChatTypePrivate)
		case "group":
			types = append(types, client.TypeChatTypeBasicGroup)
		case "supergroup":
			types = append(types, client.TypeChatTypeSupergroup, client.TypeChatTypeBasicGroup)
			if !in("group", terms) {
				justSuperGroups = true
			}
		case "secret":
			types = append(types, client.TypeChatTypeSecret)
		}
	}

	chatList, err := tgc.GetChats(chatLimit)
	if err != nil {
		slog.Error(fmt.Sprintf("Err: %s\n", err.Error()))
		os.Exit(1)
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Id", "Channel", "Group", "Supergroup", "Chat", "Bot", "Secret", "Forum", "Username", "Title"})

	idx := 0
	for _, chatId := range chatList.ChatIds {
		row := ChatRow{}

		chat, err := tgc.GetChat(chatId)
		if err != nil {
			slog.Error(fmt.Sprintf("retrieving chat %d information", chatId), "msg", err.Error())
		}

		if len(types) != 0 {
			if !in(chat.Type.ChatTypeType(), types) {
				continue
			}
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
			userInfo, err := tgc.GetUser(privateRequest.UserId)
			if err != nil {
				slog.Error("unable to get user details", "msg", err.Error())
				continue
			}

			if userInfo.Usernames != nil && len(userInfo.Usernames.ActiveUsernames) > 0 {
				row.Username = userInfo.Usernames.ActiveUsernames[0]
			}
		case client.TypeChatTypeBasicGroup:
			groupRequest := chat.Type.(*client.ChatTypeBasicGroup)
			groupInfo, err := tgc.GetGroup(groupRequest.BasicGroupId)
			if err != nil {
				slog.Error("unable to get group details", "msg", err.Error())
				continue
			}

			row.Group = true
			row.Title = chat.Title
			if len(chat.Title) > 50 {
				row.Title = chat.Title[:50] + "..."
			}

			if groupInfo.UpgradedToSupergroupId != 0 && in(client.TypeChatTypeSupergroup, types) {
				row.Supergroup = true
				superGroupInfo, err := tgc.GetSupergroup(groupInfo.UpgradedToSupergroupId)
				if err != nil {
					slog.Error("unable to get super group details", "msg", err.Error())
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
			} else if justSuperGroups {
				continue
			}
		case client.TypeChatTypeSupergroup:
			supergroupRequest := chat.Type.(*client.ChatTypeSupergroup)
			superGroupInfo, err := tgc.GetSupergroup(supergroupRequest.SupergroupId)
			if err != nil {
				slog.Error("unable to get super group details", "msg", err.Error())
				continue
			}

			row.Supergroup = true
			row.Title = chat.Title
			if len(chat.Title) > 50 {
				row.Title = chat.Title[:50] + "..."
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

			user, err := tgc.GetUser(secretRequest.UserId)
			if err != nil {
				slog.Error("unable to get user details", "msg", err.Error())
				continue
			}

			if len(user.Usernames.ActiveUsernames) > 0 {
				row.Username = user.Usernames.ActiveUsernames[0]
			}
		}

		idx++
		t.AppendRow([]interface{}{idx, row.Id, row.Channel, row.Group, row.Supergroup, row.Chat, row.Bot, row.Secret, row.Forum, row.Username, row.Title})

	}

	t.Render()
}
