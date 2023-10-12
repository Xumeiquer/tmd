package tclient

import (
	"errors"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	ApiId    int32
	ApiHash  string
	Cache    Cache
	Database Database
	Log      Log
}

type Cache struct {
	Database string
	File     string
}

type Database struct {
	Secret string
}

type Log struct {
	To    string
	Level int
}

const (
	defaultCachePath = "./.tdlib"

	apiIdNotDefined    = "API ID is not defined"
	apiHashNotDefined  = "API Hash is not defined"
	tdlibLogWrongLevel = "TDLIB log level must be between 0 and 1023"
)

var (
	ErrApiIdNotDefined   = errors.New(apiIdNotDefined)
	ErrApiHashNotDefined = errors.New(apiHashNotDefined)
	ErrTDLIBLogLevel     = errors.New(tdlibLogWrongLevel)
)

// InitConfig initializes package configuration
func InitConfig() (*Config, error) {
	var err error
	config := &Config{
		ApiId:   viper.GetInt32("tdlib.api.id"),
		ApiHash: viper.GetString("tdlib.api.hash"),
		Cache: Cache{
			Database: viper.GetString("tdlib.cache.database"),
			File:     viper.GetString("tdlib.cache.file"),
		},
		Database: Database{
			Secret: viper.GetString("tdlib.database.secret"),
		},
		Log: Log{
			To:    viper.GetString("tdlib.log.to"),
			Level: viper.GetInt("tdlib.log.level"),
		},
	}

	if config.ApiId == 0 {
		return nil, ErrApiIdNotDefined
	}
	if config.ApiHash == "" {
		return nil, ErrApiHashNotDefined
	}
	if config.Log.Level < 0 || config.Log.Level > 1023 {
		return nil, ErrTDLIBLogLevel
	}

	if config.Cache.Database == "" {
		config.Cache.Database = filepath.Join(defaultCachePath, "database")
	}
	if config.Cache.File == "" {
		config.Cache.File = filepath.Join(defaultCachePath, "files")
	}

	return config, err
}
