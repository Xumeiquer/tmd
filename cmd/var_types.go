package cmd

var (
	chatId       int64  = 0
	channelId    int64  = 0
	messageId    int64  = 0
	topicId      int64  = 0
	fileId       string = ""
	messageUrl   string = ""
	chatLimit    int32  = 1000
	showProgress bool   = false
	logLevel     string = ""
	logType      string = "json"
	logTo        string = "stdout"

	filterDate      string = initialDate
	filterContent   string = ""
	filterContentRE string = ""

	storePath string = "./media"
)

const (
	dateParserLayout = "02/01/2006"
	initialDate      = "14/08/2013"
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

type SearchResult struct {
	FileId      int32
	MessageId   int64
	AddedOn     int32
	CompletedOn int32
	IsPaused    bool
	ChatId      int64
}
