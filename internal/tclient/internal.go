package tclient

import (
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/zelenin/go-tdlib/client"
)

const (
	SystemVersion      = "v0.1.0-dev"
	ApplicationVersion = "v0.1.0-dev"
)

type TClient struct {
	cfg             *Config
	t               *client.Client
	downloadTracker map[int32]*client.File
	mux             sync.RWMutex
}

func NewTGClient() *TClient {
	config, err := InitConfig()
	if err != nil {
		slog.Error("found an error while initializing TDLib", "msg", err.Error())
	}

	return &TClient{
		cfg:             config,
		downloadTracker: make(map[int32]*client.File),
	}
}

func (tc *TClient) Stop() {
	tc.t.Stop()
}

func (tc *TClient) GetListener() *client.Listener {
	return tc.t.GetListener()
}

func (tc *TClient) GetSynchChannel() chan struct{ Id int32 } {
	return make(chan struct{ Id int32 })
}

func (tc *TClient) Authenticate() bool {
	authorizer := client.ClientAuthorizer()
	go client.CliInteractor(authorizer)

	authorizer.TdlibParameters <- &client.SetTdlibParametersRequest{
		UseTestDc:              false,
		DatabaseEncryptionKey:  []byte(tc.cfg.Database.Secret),
		DatabaseDirectory:      tc.cfg.Cache.Database,
		FilesDirectory:         tc.cfg.Cache.File,
		UseFileDatabase:        true,
		UseChatInfoDatabase:    true,
		UseMessageDatabase:     true,
		UseSecretChats:         false,
		ApiId:                  tc.cfg.ApiId,
		ApiHash:                tc.cfg.ApiHash,
		SystemLanguageCode:     "en",
		DeviceModel:            "Telegram Media Downloader",
		SystemVersion:          SystemVersion,
		ApplicationVersion:     ApplicationVersion,
		EnableStorageOptimizer: true,
		IgnoreFileNames:        false,
	}

	var err error

	_, err = client.SetLogVerbosityLevel(&client.SetLogVerbosityLevelRequest{
		NewVerbosityLevel: int32(tc.cfg.Log.Level),
	})
	if err != nil {
		slog.Error("found an error while configuring SetLogVerbosityLevel", "msg", err.Error())
		os.Exit(-1)
	}

	if tc.cfg.Log.To != "" {
		_, err = client.SetLogStream(&client.SetLogStreamRequest{
			LogStream: &client.LogStreamFile{
				Path:           tc.cfg.Log.To,
				MaxFileSize:    100000000,
				RedirectStderr: false,
			},
		})
		if err != nil {
			slog.Error("found and error while configuring SetLogStream", "msg", err.Error())
			os.Exit(-1)
		}
	}

	tc.t, err = client.NewClient(authorizer)
	if err != nil {
		slog.Error("found an error while initializing a TDLib NewClient", "msg", err.Error())
		os.Exit(-1)
	}

	return true
}

func (tc *TClient) GetUser(id int64) (*client.User, error) {
	request := &client.GetUserRequest{
		UserId: id,
	}

	return tc.t.GetUser(request)
}

func (tc *TClient) GetSecretChat(id int32) (*client.SecretChat, error) {
	request := &client.GetSecretChatRequest{
		SecretChatId: id,
	}

	return tc.t.GetSecretChat(request)
}

func (tc *TClient) GetChat(id int64) (*client.Chat, error) {
	request := client.GetChatRequest{
		ChatId: id,
	}

	return tc.t.GetChat(&request)
}

func (tc *TClient) GetGroup(id int64) (*client.BasicGroup, error) {
	request := &client.GetBasicGroupRequest{
		BasicGroupId: id,
	}

	return tc.t.GetBasicGroup(request)
}

func (tc *TClient) GetSupergroup(id int64) (*client.Supergroup, error) {
	request := &client.GetSupergroupRequest{
		SupergroupId: id,
	}

	return tc.t.GetSupergroup(request)
}

func (tc *TClient) GetMessage(chatId, messageId int64) (*client.Message, error) {
	request := &client.GetMessageRequest{
		ChatId:    chatId,
		MessageId: messageId,
	}

	return tc.t.GetMessage(request)
}

func (tc *TClient) GetMessageLinkInfo(url string) (*client.MessageLinkInfo, error) {
	request := &client.GetMessageLinkInfoRequest{
		Url: url,
	}
	return tc.t.GetMessageLinkInfo(request)
}

func (tc *TClient) GetChats(limit int32) (*client.Chats, error) {
	request := &client.GetChatsRequest{
		ChatList: nil,
		Limit:    limit,
	}

	return tc.t.GetChats(request)
}

func (tc *TClient) GetForumTopics(chatId int64, filterDate time.Time, filterContent, filterContentRE string) ([]*client.ForumTopic, error) {
	getBatch := func(offsetDate int32, offsetMessageId, offsetMessageThreadId int64) (*client.ForumTopics, error) {
		request := &client.GetForumTopicsRequest{
			ChatId:                chatId,
			Query:                 filterContent,
			OffsetDate:            offsetDate,
			OffsetMessageId:       offsetMessageId,
			OffsetMessageThreadId: offsetMessageThreadId,
			Limit:                 100,
		}

		return tc.t.GetForumTopics(request)
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
				} else if filterContentRE != "" {
					re = regexp.MustCompile(filterContentRE)
				}
				if re != nil {
					re := regexp.MustCompile(filterContentRE)

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

func (tc *TClient) GetForumTopic(chatId int64, topicId int64) (*client.ForumTopic, error) {
	request := &client.GetForumTopicRequest{
		ChatId:          chatId,
		MessageThreadId: topicId,
	}

	return tc.t.GetForumTopic(request)
}

func (tc *TClient) GetForumHistory(chatId int64, messageThreadId int64, filterDate time.Time, filterContent, filterContentRE string) ([]*client.Message, error) {
	getBatch := func(messageFromId int64) (*client.FoundChatMessages, error) {
		request := &client.SearchChatMessagesRequest{
			ChatId:          chatId,
			MessageThreadId: messageThreadId,
			FromMessageId:   messageFromId,
			Query:           filterContent,
			Offset:          0,
			Limit:           100,
		}

		return tc.t.SearchChatMessages(request)
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
					} else if filterContentRE != "" {
						re = regexp.MustCompile(filterContentRE)
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

func (tc *TClient) GetChatHistory(chatId int64, filterDate time.Time, filterContent, filterContentRE string) ([]*client.Message, error) {
	getBatch := func(messageId int64) (*client.Messages, error) {
		request := &client.GetChatHistoryRequest{
			ChatId:        chatId,
			FromMessageId: messageId,
			Offset:        0,
			Limit:         100,
			OnlyLocal:     false,
		}

		return tc.t.GetChatHistory(request)
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
					} else if filterContentRE != "" {
						re = regexp.MustCompile(filterContentRE)
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

func (tc *TClient) GetFilesFromChat(chatId int64, dateFilter time.Time, filterContent string, filterContentRE string) {
}

func (tc *TClient) DownloadFile(chatId int64, fileId string) (*client.File, error) {
	remoteRequest := &client.GetRemoteFileRequest{
		RemoteFileId: fileId,
		FileType:     nil,
	}

	remoteFile, err := tc.t.GetRemoteFile(remoteRequest)
	if err != nil {
		return nil, err
	}

	request := &client.DownloadFileRequest{
		FileId:      remoteFile.Id,
		Priority:    1,
		Offset:      0,
		Limit:       0,
		Synchronous: false,
	}

	file, err := tc.t.DownloadFile(request)

	tc.addDownload(file)

	return file, err
}
