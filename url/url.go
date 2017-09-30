package url

import (
	"fmt"
	"html"
	"net/http"
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

	if util.IgnoreCmd(cmd) {
		return "", nil
	}

	URL := ExtractURL(cmd.Raw)
	if URL == "" || util.TweetRe.MatchString(URL) || util.InstagramRe.MatchString(URL) {
		// ignore tweets and instagram posts
		return "", nil
	}
	fmt.Printf("URL:'%s'\n", URL)

	// shorten URL in goroutine
	shortChan := make(chan string)
	go func() {
		if len(URL) >= 40 {
			s, err := shortenURL(URL)
			if err == nil {
				shortChan <- s
			} else {
				shortChan <- ""
			}
		} else {
			shortChan <- ""
		}
	}()

	res, err := http.Head(URL)
	if err != nil {
		return "", err
	}

	var title string
	if res.Header["Content-Type"] == nil {
		return "", nil
	}

	if strings.Contains(res.Header["Content-Type"][0], "html") {
		// only fetch body+title for html resources
		body, err := web.GetBody(URL)
		if err != nil {
			return "", err
		}

		title = re.FindString(string(body))
		if title == "" {
			return "", nil
		}

		title = strings.Replace(title, "\n", "", -1)
		title = title[strings.Index(title, ">")+1 : strings.LastIndex(title, "<")]
		title = strings.TrimSpace(html.UnescapeString(title))

	} else {
		length := ""
		if res.Header["Content-Length"] != nil {
			length = res.Header["Content-Length"][0]
		}
		title = fmt.Sprintf("%s; Content-Length: %s", res.Header["Content-Type"][0], length)
	}

	var msg string

	shortURL := <-shortChan
	if shortURL != "" {
		msg = fmt.Sprintf("[ %s ] ", shortURL)
	}

	msg += title

	return msg, nil
}

func init() {
	bot.RegisterPassiveCommand(
		"url",
		urlTitle)
}
