/*
Copyright © 2024 Jaume Martin <jaumartin@gmail.com>
*/
package cmd

func in(val string, list []string) bool {
	for _, v := range list {
		if val == v {
			return true
		}
	}
	return false
}
