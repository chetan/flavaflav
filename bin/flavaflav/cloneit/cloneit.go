package cloneit

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/chetan/flavaflav/reddit"
	"github.com/pkg/errors"
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

type Link struct {
	Url    string
	Title  string
	Author string
}

// Copyright 2012 Jimmy Zelinskie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Submission represents an individual post from the perspective
// of a subreddit. Remember to check for nil pointers before
// using any pointer fields.
type Submission struct {
	Author        string  `json:"author"`
	Title         string  `json:"title"`
	URL           string  `json:"url"`
	Domain        string  `json:"domain"`
	Subreddit     string  `json:"subreddit"`
	SubredditID   string  `json:"subreddit_id"`
	FullID        string  `json:"name"`
	ID            string  `json:"id"`
	Permalink     string  `json:"permalink"`
	Selftext      string  `json:"selftext"`
	SelftextHTML  string  `json:"selftext_html"`
	ThumbnailURL  string  `json:"thumbnail"`
	DateCreated   float64 `json:"created_utc"`
	NumComments   int     `json:"num_comments"`
	Score         int     `json:"score"`
	Ups           int     `json:"ups"`
	Downs         int     `json:"downs"`
	IsNSFW        bool    `json:"over_18"`
	IsSelf        bool    `json:"is_self"`
	WasClicked    bool    `json:"clicked"`
	IsSaved       bool    `json:"saved"`
	BannedBy      *string `json:"banned_by"`
	LinkFlairText string  `json:"link_flair_text"`
}

var (
	pendingLinks = []*Link{}
	mtx          = sync.Mutex{}
	enabled      = false
	pluginConfig = &config{}
	errs         = 0
)

func AddLink(link *Link, channel string) {
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

func postLink(link *Link) {
	// cleanup incoming
	link.Author = cleanNick(link.Author)

	// submit
	submission, err := submit(link, pluginConfig.subreddit)
	if err != nil {
		fmt.Println("error submitting: ", err)
		return
	}

	// auto approve it
	err = approve(submission)
	if err != nil {
		fmt.Printf("failed to approve submission: %s\n%#v\n", err, submission)
		return
	}
}

func approve(submission *Submission) error {
	v := url.Values{
		"id": {submission.FullID},
	}

	type generic struct {
		Json struct {
			Errors [][]string
		}
	}

	body := strings.NewReader(v.Encode())
	err := reddit.MakeApiReq("POST", "https://oauth.reddit.com/api/approve", body, &generic{})
	if err != nil {
		return err
	}

	return nil
}

// Strip punctuation from nicks
func cleanNick(nick string) string {
	return cleanNickRE.ReplaceAllString(nick, "")
}

func submit(link *Link, subreddit string) (*Submission, error) {
	// Build form for POST request.
	v := url.Values{
		"title":       {fmt.Sprintf("<%s> %s", link.Author, link.Title)},
		"url":         {link.Url},
		"sr":          {subreddit},
		"sendreplies": {"false"},
		"resubmit":    {"false"},
		"api_type":    {"json"},
		// "flair_text":  {link.Author},
		"kind": {"link"},
	}

	type submission struct {
		Json struct {
			Errors  [][]string
			Message string `json:"message"`
			Error   string `json:"error"`
			Data    Submission
		}
	}
	submit := &submission{}

	body := strings.NewReader(v.Encode())
	err := reddit.MakeApiReq("POST", "https://oauth.reddit.com/api/submit", body, submit)
	if err != nil {
		return nil, err
	}
	if submit.Json.Error != "" {
		return nil, errors.New("failed to submit: " + submit.Json.Message)
	}

	// TODO check s.Errors and do something useful?
	return &submit.Json.Data, nil
}
