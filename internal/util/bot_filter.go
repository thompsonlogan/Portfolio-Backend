package util

import "strings"

// IsBot returns true if the User-Agent looks like a bot
func IsBot(userAgent string) bool {
	botSignatures := []string{"bot", "crawler", "spider", "slurp", "mediapartners"}
	ua := strings.ToLower(userAgent)
	for _, sig := range botSignatures {
		if strings.Contains(ua, sig) {
			return true
		}
	}
	return false
}