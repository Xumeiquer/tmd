package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

const (
	TG_API_ID           = "TMD_TG_API_ID"
	TG_API_HASH         = "TMD_TG_API_HASH"
	DB_CRYPT_PASS       = "TMD_DB_CRYPT_PASS"
	TDLIB_DB_CACHE_PATH = "TMD_TDLIB_DB_CACHE_PATH"
	TDLIB_FS_CACHE_PATH = "TMD_TDLIB_FS_CACHE_PATH"
	LOG_LEVEL           = "TMD_LOG_LEVEL"
	LOG_TO              = "TMD_LOG_TO"

	DBCryptPassValue      = ""
	TDLIBDBCachePathValue = ".tdlib/database"
	TDLIBFSCachePathValue = ".tdlib/files"

	TDLIBLogLevelMin = 0
	TDLIBLogLevelMax = 1023
)

type Config struct {
	ApiId            int32
	ApiHash          string
	DBCryptPass      string
	TDLIBDBCachePath string
	TDLIBFSCachePath string
	LogLevel         int
	LogTo            string
}

func New() *Config {
	c := &Config{}

	t, err := strconv.ParseInt(os.Getenv(TG_API_ID), 10, 32)
	if err != nil {
		slog.Error(fmt.Sprintf("Telegram API is not defined in %s", TG_API_ID))
		os.Exit(-1)
	}

	c.ApiId = int32(t)
	c.ApiHash = os.Getenv(TG_API_HASH)

	c.DBCryptPass = os.Getenv(DB_CRYPT_PASS)
	if c.DBCryptPass == "" {
		c.DBCryptPass = DBCryptPassValue
	}

	c.TDLIBDBCachePath = os.Getenv(TDLIB_DB_CACHE_PATH)
	if c.TDLIBDBCachePath == "" {
		c.TDLIBDBCachePath = TDLIBDBCachePathValue
	}

	c.TDLIBFSCachePath = os.Getenv(TDLIB_FS_CACHE_PATH)
	if c.TDLIBFSCachePath == "" {
		c.TDLIBFSCachePath = TDLIBFSCachePathValue
	}

	c.LogLevel, _ = strconv.Atoi(os.Getenv(LOG_LEVEL))
	if c.LogLevel <= 0 || c.LogLevel > 1023 {
		c.LogLevel = 1
	}

	c.LogTo = os.Getenv(LOG_TO)

	return c
}
