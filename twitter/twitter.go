package twitter

import (
	"html"
	"regexp"
	"strings"

	uri "net/url"

	"github.com/chetan/flavaflav/url"
	"github.com/go-chat-bot/bot"
	"github.com/go-chat-bot/plugins/web"
)

var (
	TweetRe   = regexp.MustCompile("https://(mobile\\.)?twitter.com/.*?/status/\\d+")
	htmlTagRe = regexp.MustCompile("<.*?>")
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

	AuthorName string `json:"author_name"`
	HTML       string
	Text       string
}

func init() {
	bot.RegisterPassiveCommand(
		"twitter",
		sniffTweet)
}

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

	if tweet.HTML != "" {
		tweet.Text = htmlTagRe.ReplaceAllString(tweet.HTML, "")
		strings.LastIndex(tweet.Text, "&mdash;")
		// tweet.Text = strings.Replace(tweet.Text, "&mdash;", " //", 1)
		tweet.Text = replaceLast(tweet.Text, "&mdash;", " //")
		tweet.Text = html.UnescapeString(tweet.Text)
	}

	return &tweet, nil
}

func sniffTweet(cmd *bot.PassiveCmd) (string, error) {
	URL := url.ExtractURL(cmd.Raw)
	if URL == "" {
		return "", nil
	}

	if TweetRe.MatchString(URL) {
		tweet, err := fetchTweet(URL)
		if err != nil {
			return "", err
		}

		return "t> " + tweet.Text, nil
	}

	return "", nil
}
