/*
Copyright © 2024 Jaume Martin <jaumartin@gmail.com>
*/
package tg

import (
	"errors"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	BasePath string
	ApiId    int32
	ApiHash  string
	Files    Files
	Database Database
	Log      Log
}

type Files struct {
	Path string
}

type Database struct {
	Secret string
	Path   string
}

type Log struct {
	To    string
	Level int
}

const (
	defaultPath = "./.tdlib"

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
		BasePath: viper.GetString("tdlib.path"),
		ApiId:    viper.GetInt32("tdlib.api.id"),
		ApiHash:  viper.GetString("tdlib.api.hash"),
		Files: Files{
			Path: viper.GetString("tdlib.files.path"),
		},
		Database: Database{
			Path:   viper.GetString("tdlib.database.path"),
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

	if config.BasePath == "" {
		config.BasePath = defaultPath
	}

	if config.Database.Path == "" {
		config.Database.Path = filepath.Join(config.BasePath, "database")
	}
	if config.Files.Path == "" {
		config.Files.Path = filepath.Join(config.BasePath, "files")
	}

	return config, err
}
