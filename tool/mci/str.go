package mci

import "strings"

// ContainsSub tests if s contains any of subs
func ContainsSub(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}

	return false
}
