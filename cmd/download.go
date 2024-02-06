/*
Copyright © 2024 Jaume Martin <jaumartin@gmail.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download media from the Telegram cloud",
	Long:  `Download command allows you to download media from messages, channels, chats, topics, etc.`,
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.PersistentFlags().StringVarP(&storePath, "store", "", storePath, "path where media will be saved")
	downloadCmd.PersistentFlags().StringVarP(&filterDate, "date", "", filterDate, "filter results by date")
	downloadCmd.PersistentFlags().StringVarP(&filterContent, "filter", "", filterContent, "filter results by content")
}
