package progress_bar

import "github.com/spf13/viper"

type Config struct {
	Log Log
}

type Log struct {
	Level string
	Type  string
	To    string
}

func Initconfig() (*Config, error) {
	var err error
	config := &Config{
		Log: Log{
			Level: viper.GetString("tmd.log.level"),
			Type:  viper.GetString("tmd.log.type"),
			To:    viper.GetString("tmd.log.to"),
		},
	}
	return config, err
}
