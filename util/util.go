package util

import "regexp"

var (
	TweetRe     = regexp.MustCompile("https://(mobile\\.)?twitter.com/.*?/status/\\d+")
	InstagramRe = regexp.MustCompile(`^https://www.instagram.com/p/.*?/`)
)
