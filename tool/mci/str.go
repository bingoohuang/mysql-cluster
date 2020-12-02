package mci

import "strings"

// ContainsSub tests if s contains any of subs.
func ContainsSub(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}

	return false
}

// EqualsFold tests if s equals fold any of subs.
func EqualsFold(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.EqualFold(s, sub) {
			return true
		}
	}

	return false
}
