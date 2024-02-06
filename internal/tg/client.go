/*
Copyright © 2024 Jaume Martin <jaumartin@gmail.com>
*/
package tg

import (
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"time"

	"github.com/zelenin/go-tdlib/client"
)

const (
	SystemVersion      = "v0.1.0"
	ApplicationVersion = "v0.1.0"
)

type TGClient struct {
	cfg          *Config
	c            *client.Client
	storePath    string
	showProgress bool
	dm           DownloadManager
}

func NewTGClient(options ...func(*TGClient)) *TGClient {
	config, err := InitConfig()
	if err != nil {
		slog.Error("found an error while initializing TDLib", "msg", err.Error())
	}

	tgcli := &TGClient{
		cfg: config,
	}

	for _, o := range options {
		o(tgcli)
	}

	return tgcli
}

func (tg *TGClient) Stop() {
	tg.c.Stop()
}

func (tg *TGClient) Authenticate() bool {
	authorizer := client.ClientAuthorizer()
	go client.CliInteractor(authorizer)

	authorizer.TdlibParameters <- &client.SetTdlibParametersRequest{
		UseTestDc:              false,
		DatabaseEncryptionKey:  []byte(tg.cfg.Database.Secret),
		DatabaseDirectory:      tg.cfg.Database.Path,
		FilesDirectory:         tg.cfg.Files.Path,
		UseFileDatabase:        true,
		UseChatInfoDatabase:    true,
		UseMessageDatabase:     true,
		UseSecretChats:         false,
		ApiId:                  tg.cfg.ApiId,
		ApiHash:                tg.cfg.ApiHash,
		SystemLanguageCode:     "en",
		DeviceModel:            "Telegram Media Downloader",
		SystemVersion:          SystemVersion,
		ApplicationVersion:     ApplicationVersion,
		EnableStorageOptimizer: true,
		IgnoreFileNames:        false,
	}

	var err error

	_, err = client.SetLogVerbosityLevel(&client.SetLogVerbosityLevelRequest{
		NewVerbosityLevel: int32(tg.cfg.Log.Level),
	})
	if err != nil {
		slog.Error("found an error while configuring SetLogVerbosityLevel", "msg", err.Error())
		os.Exit(-1)
	}

	if tg.cfg.Log.To != "" {
		_, err = client.SetLogStream(&client.SetLogStreamRequest{
			LogStream: &client.LogStreamFile{
				Path:           tg.cfg.Log.To,
				MaxFileSize:    100000000,
				RedirectStderr: false,
			},
		})
		if err != nil {
			slog.Error("found and error while configuring SetLogStream", "msg", err.Error())
			os.Exit(-1)
		}
	}

	tg.c, err = client.NewClient(authorizer)
	if err != nil {
		slog.Error("found an error while initializing a TDLib NewClient", "msg", err.Error())
		os.Exit(-1)
	}

	return true
}

func WithStorePath(storePath string) func(*TGClient) {
	return func(tgcli *TGClient) {
		tgcli.storePath = storePath
	}
}

func WithPrintStatus(showProgress bool) func(*TGClient) {
	return func(tgcli *TGClient) {
		tgcli.showProgress = showProgress
	}
}

func (tg *TGClient) SetDownloadManager(dm DownloadManager) {
	tg.dm = dm
}

func (tg *TGClient) GetListener() *client.Listener {
	return tg.c.GetListener()
}

func (tg *TGClient) GetChat(id int64) (*client.Chat, error) {
	request := client.GetChatRequest{
		ChatId: id,
	}

	return tg.c.GetChat(&request)
}

func (tg *TGClient) GetChats(limit int32) (*client.Chats, error) {
	request := &client.GetChatsRequest{
		ChatList: nil,
		Limit:    limit,
	}

	return tg.c.GetChats(request)
}

func (tg *TGClient) GetUser(id int64) (*client.User, error) {
	request := &client.GetUserRequest{
		UserId: id,
	}

	return tg.c.GetUser(request)
}

func (tg *TGClient) GetGroup(id int64) (*client.BasicGroup, error) {
	request := &client.GetBasicGroupRequest{
		BasicGroupId: id,
	}

	return tg.c.GetBasicGroup(request)
}

func (tg *TGClient) GetSupergroup(id int64) (*client.Supergroup, error) {
	request := &client.GetSupergroupRequest{
		SupergroupId: id,
	}

	return tg.c.GetSupergroup(request)
}

func (tg *TGClient) GetForumTopics(chatId int64, filterDate time.Time, filterContent string) ([]*client.ForumTopic, error) {
	getBatch := func(offsetDate int32, offsetMessageId, offsetMessageThreadId int64) (*client.ForumTopics, error) {
		request := &client.GetForumTopicsRequest{
			ChatId:                chatId,
			Query:                 filterContent,
			OffsetDate:            offsetDate,
			OffsetMessageId:       offsetMessageId,
			OffsetMessageThreadId: offsetMessageThreadId,
			Limit:                 100,
		}

		return tg.c.GetForumTopics(request)
	}

	topics := []*client.ForumTopic{}
	filterOffsetDate := int32(0)
	lastOffsetMessageId := int64(0)
	lastOffsetMessageThreadId := int64(0)

	filterDateInt := int32(filterDate.Unix())

	end := false

	for !end {
		lastTopics, err := getBatch(filterOffsetDate, lastOffsetMessageId, lastOffsetMessageThreadId)
		if err != nil {
			return topics, err
		}

		if len(lastTopics.Topics) == 0 {
			break
		}

		filterOffsetDate = lastTopics.NextOffsetDate
		lastOffsetMessageId = lastTopics.NextOffsetMessageId
		lastOffsetMessageThreadId = lastTopics.NextOffsetMessageThreadId

		for _, lastTopic := range lastTopics.Topics {
			if lastTopic.Info.CreationDate >= filterDateInt {
				var re *regexp.Regexp
				if filterContent != "" {
					re = regexp.MustCompile(fmt.Sprintf(".*%s.*", filterContent))
				} else if filterContent != "" {
					re = regexp.MustCompile(filterContent)
				}
				if re != nil {
					re := regexp.MustCompile(filterContent)

					if regexFilter(re, lastTopic.Info.Name) {
						topics = append(topics, lastTopic)
					}
				} else {
					topics = append(topics, lastTopic)
				}
			}
		}
	}

	return topics, nil
}

func (tg *TGClient) GetMessageLinkInfo(url string) (*client.MessageLinkInfo, error) {
	request := &client.GetMessageLinkInfoRequest{
		Url: url,
	}
	return tg.c.GetMessageLinkInfo(request)
}

func (tg *TGClient) GetForumHistory(chatId int64, messageThreadId int64, filterDate time.Time, filterContent string) ([]*client.Message, error) {
	getBatch := func(messageFromId int64) (*client.FoundChatMessages, error) {
		request := &client.SearchChatMessagesRequest{
			ChatId:          chatId,
			MessageThreadId: messageThreadId,
			FromMessageId:   messageFromId,
			Query:           filterContent,
			Offset:          0,
			Limit:           100,
		}

		return tg.c.SearchChatMessages(request)
	}

	messages := []*client.Message{}

	filterDateInt := int32(filterDate.Unix())
	lastMessageFromId := int64(0)
	end := false

	for !end {
		lastMessages, err := getBatch(lastMessageFromId)
		if err != nil {
			return messages, err
		}
		if len(lastMessages.Messages) == 0 {
			break
		}

		lastMessageFromId = lastMessages.NextFromMessageId

		for _, lastMessage := range lastMessages.Messages {
			if lastMessage.Date >= filterDateInt {
				switch lastMessage.Content.MessageContentType() {
				case client.TypeMessageDocument:
					var re *regexp.Regexp
					if filterContent != "" {
						re = regexp.MustCompile(fmt.Sprintf(".*%s.*", filterContent))
					}
					if re != nil {
						if content, ok := lastMessage.Content.(*client.MessageDocument); ok {
							if regexFilter(re, content.Caption.Text, content.Document.FileName, content.Document.Type) {
								messages = append(messages, lastMessage)
							}
						}
					} else {
						messages = append(messages, lastMessage)
					}
				}
			} else {
				end = true
			}
		}
	}

	return messages, nil
}

func (tg *TGClient) GetChatHistory(chatId int64, filterDate time.Time, filterContent string) ([]*client.Message, error) {
	getBatch := func(messageId int64) (*client.Messages, error) {
		request := &client.GetChatHistoryRequest{
			ChatId:        chatId,
			FromMessageId: messageId,
			Offset:        0,
			Limit:         100,
			OnlyLocal:     false,
		}

		return tg.c.GetChatHistory(request)
	}

	messages := []*client.Message{}

	filterDateInt := int32(filterDate.Unix())
	lastMessageId := int64(0)
	end := false

	for !end {
		lastMessages, err := getBatch(lastMessageId)
		if err != nil {
			return nil, err
		}

		if len(lastMessages.Messages) == 0 {
			break
		}

		for _, lastMessage := range lastMessages.Messages {
			if lastMessage.Date >= filterDateInt {
				switch lastMessage.Content.MessageContentType() {
				case client.TypeMessageDocument:
					var re *regexp.Regexp
					if filterContent != "" {
						re = regexp.MustCompile(fmt.Sprintf(".*%s.*", filterContent))
					}
					if re != nil {
						if content, ok := lastMessage.Content.(*client.MessageDocument); ok {
							if regexFilter(re, content.Caption.Text, content.Document.FileName, content.Document.Type) {
								messages = append(messages, lastMessage)
							}
						}
					} else {
						messages = append(messages, lastMessage)
					}
				}
				lastMessageId = lastMessage.Id
			} else {
				end = true
			}
		}
	}

	return messages, nil
}

func (tg *TGClient) DownloadFile(chatId int64, fileId string) (*client.File, error) {
	remoteRequest := &client.GetRemoteFileRequest{
		RemoteFileId: fileId,
		FileType:     nil,
	}

	remoteFile, err := tg.c.GetRemoteFile(remoteRequest)
	if err != nil {
		return nil, err
	}

	downloadRequest := &client.DownloadFileRequest{
		FileId:      remoteFile.Id,
		Priority:    1,
		Offset:      0,
		Limit:       0,
		Synchronous: false,
	}

	file, err := tg.c.DownloadFile(downloadRequest)
	if err != nil {
		return nil, err
	}

	tg.dm.AddFile(file, tg.storePath)

	return file, err
}
