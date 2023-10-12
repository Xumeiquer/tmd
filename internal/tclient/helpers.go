package tclient

import "regexp"

func regexFilter(re *regexp.Regexp, args ...string) bool {
	for _, arg := range args {
		if re.FindAllString(arg, -1) != nil {
			return true
		}
	}
	return false
}
