package url

import (
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strings"

	"github.com/chetan/flavaflav/util"
	"github.com/go-chat-bot/bot"
	"github.com/go-chat-bot/plugins/web"
)

const (
	minDomainLength = 3
)

var (
	re = regexp.MustCompile("<title( .*?)?>\\n*?(.*?)\\n*?<\\/title>")
)

func canBeURLWithoutProtocol(text string) bool {
	return len(text) > minDomainLength &&
		!strings.HasPrefix(text, "http") &&
		strings.Contains(text, ".")
}

func ExtractURL(text string) string {
	extractedURL := ""
	for _, value := range strings.Split(text, " ") {
		if canBeURLWithoutProtocol(value) {
			value = "http://" + value
		}

		parsedURL, err := url.Parse(value)
		if err != nil {
			continue
		}
		if strings.HasPrefix(parsedURL.Scheme, "http") {
			extractedURL = parsedURL.String()
			break
		}
	}
	return extractedURL
}

func shortenURL(u string) (string, error) {
	return shortenURL99c(u)
}

func urlTitle(cmd *bot.PassiveCmd) (string, error) {
	URL := ExtractURL(cmd.Raw)
	if URL == "" || util.TweetRe.MatchString(URL) || util.InstagramRe.MatchString(URL) {
		// ignore tweets and instagram posts
		return "", nil
	}
	fmt.Printf("URL:'%s'\n", URL)

	body, err := web.GetBody(URL)
	if err != nil {
		return "", err
	}

	title := re.FindString(string(body))
	if title == "" {
		return "", nil
	}

	title = strings.Replace(title, "\n", "", -1)
	title = title[strings.Index(title, ">")+1 : strings.LastIndex(title, "<")]
	title = strings.TrimSpace(html.UnescapeString(title))

	var msg string

	if len(URL) >= 40 {
		// shorten URLs longer than 40
		shortURL, err := shortenURL(URL)
		if err == nil {
			msg = fmt.Sprintf("[ %s ] ", shortURL)
		}
	}

	msg = fmt.Sprintf("%s%s", msg, title)

	return msg, nil
}

func init() {
	bot.RegisterPassiveCommand(
		"url",
		urlTitle)
}
