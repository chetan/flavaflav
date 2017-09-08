package twitter

import (
	"fmt"
	"html"
	"regexp"
	"strings"

	uri "net/url"

	"github.com/chetan/flavaflav/url"
	"github.com/chetan/flavaflav/util"
	"github.com/go-chat-bot/bot"
	"github.com/go-chat-bot/plugins/web"
)

var (
	htmlTagRe = regexp.MustCompile("<.*?>")
	handleRe  = regexp.MustCompile(`\((@.*?)\) (.*)`)
)

type Tweet struct {
	CacheAge     string `json:"cache_age"`
	AuthorURL    string `json:"author_url"`
	ProviderName string `json:"provider_name"`
	ProviderURL  string `json:"provider_url"`
	URL          string
	Type         string
	Width        int
	Height       int
	Version      string

	AuthorName   string `json:"author_name"`
	HTML         string
	Text         string
	TextBody     string
	TextDate     string
	AuthorHandle string
}

func init() {
	bot.RegisterPassiveCommand(
		"twitter",
		sniffTweet)
}

// replaceLast replaces the last occurence of the needle with the given string
func replaceLast(str string, needle string, replace string) string {
	i := strings.LastIndex(str, needle)
	return str[0:i] + replace + str[i+len(needle):]
}

func fetchTweet(u string) (*Tweet, error) {
	tweet := Tweet{}
	err := web.GetJSON("https://publish.twitter.com/oembed?url="+uri.QueryEscape(u), &tweet)
	if err != nil {
		return nil, err
	}
	processTweet(&tweet)
	return &tweet, nil
}

func processTweet(tweet *Tweet) {
	if tweet.HTML != "" {
		tweet.Text = htmlTagRe.ReplaceAllString(tweet.HTML, "")

		// extract body
		i := strings.LastIndex(tweet.Text, "&mdash;")
		tweet.TextBody = html.UnescapeString(strings.TrimSpace(tweet.Text[0:i]))

		// extract handle
		matches := handleRe.FindStringSubmatch(tweet.Text[i+7:])
		tweet.AuthorHandle = matches[1]
		tweet.TextDate = matches[2]

		// tweet.Text = strings.Replace(tweet.Text, "&mdash;", " //", 1)
		tweet.Text = replaceLast(tweet.Text, "&mdash;", " //")
		tweet.Text = html.UnescapeString(tweet.Text)
	}
}

func sniffTweet(cmd *bot.PassiveCmd) (string, error) {
	URL := url.ExtractURL(cmd.Raw)
	if URL == "" {
		return "", nil
	}

	if util.TweetRe.MatchString(URL) {
		tweet, err := fetchTweet(URL)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("<%s> %s // %s", tweet.AuthorHandle, tweet.TextBody, tweet.TextDate), nil
	}

	return "", nil
}
