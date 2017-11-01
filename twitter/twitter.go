package twitter

import (
	"fmt"
	"html"
	"regexp"
	"strings"

	uri "net/url"

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

func (t Tweet) String() string {
	return fmt.Sprintf("<%s> %s // %s", t.AuthorHandle, t.TextBody, t.TextDate)
}

func init() {
	bot.RegisterPassiveCommand(
		"twitter",
		handleTweet)
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

func handleTweet(cmd *bot.PassiveCmd) (string, error) {
	if util.IgnoreCmd(cmd) {
		return "", nil
	}

	URL := util.ExtractURL(cmd.Raw)
	if URL == "" {
		return "", nil
	}

	if util.IsTwitter(URL) {
		tweet, err := fetchTweet(URL)
		if err != nil {
			return "", err
		}

		out := tweet.String()

		/* // skip embedded stuff for now
		out := Gray(tweet.String())
		embeddedURLs := util.ExtractURLs(tweet.TextBody)
		for _, u := range embeddedURLs {
			fmt.Println(u)
			if util.IsTwitterShortURL(u) {
				expanded, err := util.ExpandURL(u)
				if err == nil && expanded != "" && util.IsTwitter(expanded) {
					t, err := fetchTweet(expanded)
					if err == nil {
						out += "\n" + Gray(" \\--- "+t.String())
					}
				}
			}
		}
		*/

		return out, nil
	}

	return "", nil
}

const (
	WHITE = iota
	BLACK
	NAVY
	GREEN
	RED
	MAROON
	PURPLE
	OLIVE
	YELLOW
	LIGHTGREEN
	TEAL
	CYAN
	ROYALBLUE
	MAGENTA
	GRAY
	LIGHTGRAY
)

func Color(s string, c int) string {
	return fmt.Sprintf("\x03%d%s\x03", c, s)
}

func Gray(s string) string {
	return Color(s, GRAY)
}
