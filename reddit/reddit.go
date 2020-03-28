package reddit

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// Link for submitting to reddit
type Link struct {
	Url    string
	Title  string
	Author string
}

// Submission response type, courtesy of:
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

// Submit the given link to a subreddit
func Submit(link *Link, subreddit string) (*Submission, error) {
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
	err := MakeApiReq("POST", "https://oauth.reddit.com/api/submit", body, submit)
	if err != nil {
		return nil, err
	}
	if submit.Json.Error != "" {
		return nil, errors.New("failed to submit: " + submit.Json.Message)
	}

	// TODO check s.Errors and do something useful?
	return &submit.Json.Data, nil
}

// Approve the given submission
func Approve(submission *Submission) error {
	v := url.Values{
		"id": {submission.FullID},
	}

	type generic struct {
		Json struct {
			Errors [][]string
		}
	}

	body := strings.NewReader(v.Encode())
	err := MakeApiReq("POST", "https://oauth.reddit.com/api/approve", body, &generic{})
	if err != nil {
		return err
	}

	return nil
}
