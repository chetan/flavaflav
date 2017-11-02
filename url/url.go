package url

import (
	"fmt"
	"html"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"code.cloudfoundry.org/bytefmt"

	"github.com/chetan/flavaflav/util"
	"github.com/go-chat-bot/bot"
	"github.com/go-chat-bot/plugins/web"
)

var (
	titleRe = regexp.MustCompile(`<title( .*?)?>\n*?(.*?)\n*?<\\?/title>`)
)

func shortenURL(u string) (string, error) {
	return shortenURL99c(u)
}

func urlTitle(cmd *bot.PassiveCmd) (string, error) {

	if util.IgnoreCmd(cmd) {
		return "", nil
	}

	URL := util.ExtractURL(cmd.Raw)
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

	// send a HEAD request first to grab content-type & length headers
	req, err := http.NewRequest("HEAD", URL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	var title string
	if res.Header["Content-Type"] == nil {
		return "", nil
	}

	if strings.Contains(res.Header["Content-Type"][0], "html") {
		title, err = extractTitle(URL)
		if err != nil {
			return "", err
		}

	} else {
		length := ""
		if res.Header["Content-Length"] != nil {
			length = res.Header["Content-Length"][0]
			if length != "" {
				i, err := strconv.ParseUint(length, 10, 64)
				if err == nil {
					length = bytefmt.ByteSize(i)
				}
			}
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

func extractTitle(url string) (string, error) {
	// only fetch body+title for html resources
	body, err := web.GetBody(url)
	if err != nil {
		return "", err
	}

	var title string

	// workaround for jalopnik.com properties - they return a whole mess of JS
	// with embedded HTML on a single line before the actual <title> tag. The real
	// <title> is the last such match.
	titles := titleRe.FindAllString(string(body), -1)

	if len(titles) > 0 {
		title = titles[len(titles)-1]
	}

	if title == "" {
		return "", nil
	}

	title = strings.Replace(title, "\n", "", -1)
	title = title[strings.Index(title, ">")+1 : strings.LastIndex(title, "<")]
	title = strings.TrimSpace(html.UnescapeString(title))
	return title, nil
}

func init() {
	bot.RegisterPassiveCommand(
		"url",
		urlTitle)
}
