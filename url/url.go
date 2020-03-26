package url

import (
	"fmt"
	"html"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"code.cloudfoundry.org/bytefmt"

	"github.com/chetan/flavaflav/bin/flavaflav/cloneit"
	"github.com/chetan/flavaflav/util"
	"github.com/go-chat-bot/bot"
)

var (
	titleRe = regexp.MustCompile(`<title( .*?)?>\n*?(.*?)\n*?<\\?/title>`)
)

func ShortenURL(u string) (string, error) {
	return shortenURL99c(u)
}

func urlTitle(cmd *bot.PassiveCmd) (string, error) {

	if util.IgnoreCmd(cmd) {
		return "", nil
	}

	URL := util.ExtractURL(cmd.Raw)
	if URL == "" || util.IsTwitter(URL) || util.IsInstagram(URL) || util.IsYoutube(URL) {
		// ignore tweets and instagram posts
		// ignore youtube due to other bot
		return "", nil
	}
	fmt.Printf("URL:'%s'\n", URL)

	// shorten URL in goroutine
	shortChan := make(chan string)
	go func() {
		if len(URL) >= 40 {
			s, err := ShortenURL(URL)
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
	util.AddHeaders(req)
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

	// Send to cloneit
	cloneit.AddLink(&cloneit.Link{
		Url:     URL,
		Title:   title,
		Channel: cmd.Channel,
		Author:  cmd.User.Nick,
	})

	msg += title

	return msg, nil
}

func extractTitle(url string) (string, error) {
	// only fetch body+title for html resources
	body, err := util.GetBody(url)
	if err != nil {
		return "", err
	}

	var title string

	// workaround for kinja.com properties - they return a whole mess of JS
	// with embedded HTML on a single line before the actual <title> tag. The real
	// <title> is the last such match.
	if util.IsKinjaNetwork(url) {
		titles := titleRe.FindAllString(string(body), -1)
		if len(titles) > 0 {
			// title comes at the end
			title = titles[len(titles)-1]
		}
	} else {
		title = titleRe.FindString(string(body))
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
