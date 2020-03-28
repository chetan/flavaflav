package cloneit

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/chetan/flavaflav/reddit"
)

var cleanNickRE = regexp.MustCompile(`[-_^]`)

type config struct {
	clientID     string
	clientSecret string
	username     string
	password     string
	subreddit    string
	channel      string
	access       string
	refresh      string
}

var (
	enabled      = false
	pluginConfig = &config{}
)

func AddLink(link *reddit.Link, channel string) {
	if enabled && channel == pluginConfig.channel {
		go postLink(link) // submit in background
	}
}

// Enable the plugin
func Enable(clientID, clientSecret, user string, pass string, sub string, channel string, access string, refresh string) {
	pluginConfig = &config{clientID, clientSecret, user, pass, sub, channel, access, refresh}

	reddit.SetClientID(clientID)
	reddit.SetClientSecret(clientSecret)
	reddit.SetUserAgent("flavaflav v1")

	if refresh == "" {
		fmt.Println("enabling reddit plugin for the first time. trying to oauth handshake....")
		scopes := strings.Join([]string{"identity", "flair", "modflair", "modposts", "mysubreddits", "submit"}, ",")
		a, r, err := reddit.PerformHandshake("http://localhost", []string{scopes}, true)
		if err != nil {
			panic("failed to handshake with reddit")
		}
		fmt.Println("reddit handshake success!")
		fmt.Println("reddit_access_token:", a)
		fmt.Println("reddit_refresh_token:", r)
		fmt.Println("Store these tokens in ~/.flavaflav.yml")
		os.Exit(0)
	}

	reddit.SetAccessToken(access)
	reddit.SetRefreshToken(refresh)

	err := reddit.RefreshCreds()
	if err != nil {
		fmt.Println("error refreshing creds: ", err)
		os.Exit(1)
	}

	enabled = true
}

func postLink(link *reddit.Link) {
	// cleanup incoming
	link.Author = cleanNick(link.Author)

	// submit
	submission, err := reddit.Submit(link, pluginConfig.subreddit)
	if err != nil {
		fmt.Println("error submitting: ", err)
		return
	}

	// auto approve it
	err = reddit.Approve(submission)
	if err != nil {
		fmt.Printf("failed to approve submission: %s\n%#v\n", err, submission)
		return
	}
}

// Strip punctuation from nicks
func cleanNick(nick string) string {
	return cleanNickRE.ReplaceAllString(nick, "")
}
