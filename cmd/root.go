/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Version       = "v0.1.0-dev"
	CommitHash    = ""
	BuildTime     = ""
	versionString = fmt.Sprintf(`Telegram Media Downloader
Version: %s
Git commit: %s
Build date: %s
`, Version, CommitHash, BuildTime)

	// configFile configuration file
	configFile string
)

const (
	tmdConfigFileName = "tmd"
	configPrefix      = "TMD"
)

// RootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "tmd",
	Version: Version,
	Short:   "Telegram Media Downloader.",
	Long: `Download Telegram media from Users, Chats, or Channels.

Telegram Media Downloader allow users to download media content from Telegram cloud 
without manually interacting with the Telegram client. Telegram Media Downloader is
or acts as a client so it has to be enrolled as a Telegram client.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.SetVersionTemplate(versionString)
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "", logLevel, "set log level (info, error, warn, debug)")
	rootCmd.PersistentFlags().StringVarP(&logType, "log-type", "", logType, "log as text or JSON")
	rootCmd.PersistentFlags().StringVarP(&logTo, "log-to", "", logTo, "where to log")

	// Bind flags to configuration file
	viper.BindPFlag("tmd.log.to", rootCmd.PersistentFlags().Lookup("log-to"))
	viper.BindPFlag("tmd.log.level", rootCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("tmd.log.type", rootCmd.PersistentFlags().Lookup("log-type"))
}

func initConfig() {
	// Define default configuration file
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName(tmdConfigFileName)
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")

		switch runtime.GOOS {
		case "windows":
			viper.AddConfigPath(filepath.Join(os.Getenv("APPDATA"), tmdConfigFileName))
		case "darwin":
			fallthrough
		case "linux":
			viper.AddConfigPath(filepath.Join("etc", tmdConfigFileName))
		}
	}

	// Read configuration from environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix(configPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			slog.Warn("configuration file not found")
		} else {
			// Config file was found but another error was produced
			slog.Error("found an error while loading the configuration file", "msg", err.Error())
		}
	}
}
