package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/go-chat-bot/bot"
	"github.com/ryanuber/go-glob"
)

var (
	VoxRe           = regexp.MustCompile(`^https?://(www.)?vox.com/`)
	TweetRe         = regexp.MustCompile("https://(mobile\\.)?twitter.com/.*?/status/\\d+")
	InstagramRe     = regexp.MustCompile(`^https://www.instagram.com/p/.*?/`)
	TwitterShortUrl = regexp.MustCompile(`https://t\.co/.*`)
	KinjaRe         = regexp.MustCompile(strings.TrimSpace(`
^https?://((.*?)\.)?(kotaku.com|lifehacker.com|splinternews.com|earther.com|deadspin.com|gizmodo.com|jalopnik.com|jezebel.com|theroot.com|kinja.com|theonion.com|avclub.com)
`))
)

const (
	// for url extractor
	minDomainLength = 3
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

func IsKinjaNetwork(u string) bool {
	return KinjaRe.MatchString(u)
}

func IsVox(u string) bool {
	return VoxRe.MatchString(u)
}

func IsTwitter(u string) bool {
	return TweetRe.MatchString(u)
}

func IsTwitterShortURL(u string) bool {
	return TwitterShortUrl.MatchString(u)
}

func IsInstagram(u string) bool {
	return InstagramRe.MatchString(u)
}

func ExpandURL(u string) (string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(u)
	if err != nil && err != http.ErrUseLastResponse {
		return "", err
	}

	loc := resp.Header["Location"]
	return loc[0], nil
}

func canBeURLWithoutProtocol(text string) bool {
	return len(text) > minDomainLength &&
		!strings.HasPrefix(text, "http") &&
		strings.Contains(text, ".")
}

func AddHeaders(req *http.Request) {
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.62 Safari/537.36")
}

func ExtractURLs(text string) []string {

	var urls []string

	for _, value := range strings.Split(text, " ") {
		if canBeURLWithoutProtocol(value) {
			value = "http://" + value
		}

		parsedURL, err := url.Parse(value)
		if err != nil {
			continue
		}
		if strings.HasPrefix(parsedURL.Scheme, "http") {
			urls = append(urls, parsedURL.String())
		}
	}

	return urls
}

func ExtractURL(text string) string {
	urls := ExtractURLs(text)
	if urls == nil || len(urls) == 0 {
		return ""
	}
	return urls[0]
}

func GetBody(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	AddHeaders(req)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

func GetJSON(url string, v interface{}) error {
	body, err := GetBody(url)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}
