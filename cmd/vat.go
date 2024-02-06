/*
Copyright © 2024 Jaume Martin <jaumartin@gmail.com>
*/
package cmd

const (
	dateParserLayout = "02/01/2006"
	initialDate      = "14/08/2013"
)

var (
	filterType    string
	filterContent string
	filterDate    string
	chatLimit     int32 = 1000
	channelId     int64 = 0
	chatId        int64 = 0
	topicId       int64 = 0
	fileId        string
	messageUrl    string
	storePath     string = "./media"
)

type ChatRow struct {
	Id         int64
	Channel    bool
	Group      bool
	Supergroup bool
	Chat       bool
	Bot        bool
	Secret     bool
	Forum      bool
	Username   string
	Title      string
}

type TopicRow struct {
	Id           int64
	Name         string
	General      bool
	Closed       bool
	Hidden       bool
	CreationDate int32
}

type DocumentRow struct {
	Id        string
	ChatId    int64
	MessageId int64
	Message   string
	Name      string
	Mime      string
	Size      int64
	Date      int32
}
