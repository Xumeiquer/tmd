/*
Copyright © 2024 Jaume Martin <jaumartin@gmail.com>
*/
package tg

import (
	"regexp"

	"github.com/zelenin/go-tdlib/client"
)

func regexFilter(re *regexp.Regexp, args ...string) bool {
	for _, arg := range args {
		if re.FindAllString(arg, -1) != nil {
			return true
		}
	}
	return false
}

func inDownloadQueue(file *client.File, files []*client.File) bool {
	for _, f := range files {
		if f.Id == file.Id {
			return true
		}
	}
	return false
}
