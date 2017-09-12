package util

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-chat-bot/bot"
	"github.com/ryanuber/go-glob"
)

var (
	TweetRe     = regexp.MustCompile("https://(mobile\\.)?twitter.com/.*?/status/\\d+")
	InstagramRe = regexp.MustCompile(`^https://www.instagram.com/p/.*?/`)
)

// IgnoreNicks is a list of nicks to ignore
var IgnoreNicks []string

// IgnorePatterns is a list of user patterns of the format 'nick!ident@host' to ignore
var IgnorePatterns []string

// IgnoreCmd returns true if the given command should be ignored
func IgnoreCmd(cmd *bot.PassiveCmd) bool {

	// test IgnoreUsers list
	for _, n := range IgnoreNicks {
		if n == cmd.User.Nick {
			return true
		}
	}

	// test IgnorePatterns list
	for _, p := range IgnorePatterns {
		if userMatchesPattern(cmd.User, p) {
			return true
		}
	}

	return false
}

func userMatchesPattern(user *bot.User, p string) (match bool) {
	if p == "" {
		return false
	}

	p = strings.ToLower(p)

	str := fmt.Sprintf("%s!%s@%s", user.Nick, user.RealName, user.ID)
	return glob.Glob(p, str)
}
