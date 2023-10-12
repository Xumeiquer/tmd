package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download media from the Telegram cloud",
	Long:  `Download command allows you to download media from messages, channels, chats, topics, etc.`,
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.PersistentFlags().StringVarP(&storePath, "store-path", "", storePath, "path where media will be saved")
	downloadCmd.PersistentFlags().BoolVarP(&showProgress, "show-progress", "", showProgress, "show download progress (not recommended for large downloads)")
	downloadCmd.PersistentFlags().StringVarP(&filterDate, "filter-date", "", filterDate, "filter results by date")
	downloadCmd.PersistentFlags().StringVarP(&filterContent, "filter", "", filterContent, "filter results by content")
	downloadCmd.PersistentFlags().StringVarP(&filterContentRE, "filter-re", "", filterContentRE, "filter results by regular expresion")

	// Bind flags to configuration file
	viper.BindPFlag("tmd.cmds.filters.date", downloadCmd.PersistentFlags().Lookup("filter-date"))
	viper.BindPFlag("tmd.cmds.filters.regex", downloadCmd.PersistentFlags().Lookup("filter-re"))
	viper.BindPFlag("tmd.cmds.filters.text", downloadCmd.PersistentFlags().Lookup("filter"))

	viper.BindPFlag("tmd.cmds.download.store-path", downloadCmd.PersistentFlags().Lookup("store-path"))
	viper.BindPFlag("tmd.cmds.download.show-progress", downloadCmd.PersistentFlags().Lookup("show-progress"))
}
