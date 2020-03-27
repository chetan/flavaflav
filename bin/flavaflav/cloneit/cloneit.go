package cloneit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/jzelinskie/geddit"
	"github.com/pkg/errors"
)

type config struct {
	clientId     string
	clientSecret string
	username     string
	password     string
	subreddit    string
	channel      string
}

type Link struct {
	Url     string
	Title   string
	Author  string
	Channel string
}

var (
	pendingLinks = []*Link{}
	mtx          = sync.Mutex{}
	enabled      = false
	pluginConfig = &config{}
	errs         = 0
)

var session *geddit.OAuthSession

func AddLink(link *Link) {
	if enabled && link.Channel == pluginConfig.channel {
		go postLink(link) // submit in background
	}
}

// Enable the plugin
func Enable(clientId, clientSecret, user string, pass string, sub string, channel string) {
	pluginConfig = &config{clientId, clientSecret, user, pass, sub, channel}
	sess, err := createSession()
	if err != nil {
		fmt.Println("cloneit: err creating initial session: ", err)
		return
	}
	session = sess
	enabled = true
}

func createSession() (*geddit.OAuthSession, error) {
	if errs > 2 {
		// give up
		return nil, errors.New("too many errors, giving up creating session")
	}

	sess, err := geddit.NewOAuthSession(
		pluginConfig.clientId,
		pluginConfig.clientSecret,
		"flavaflav v1",
		"http://127.0.0.1/",
	)

	err = sess.LoginAuth(pluginConfig.username, pluginConfig.password)
	if err != nil {
		fmt.Println("error logging in:", err)
	}

	if err != nil {
		// retry until limit
		fmt.Println("error creating session:", err)
		errs++
		return createSession()
	}

	errs = 0
	return sess, nil
}

func postLink(link *Link) {
	if session == nil {
		fmt.Println("reddit session is nil, not posting")
		return
	}

	if errs > 2 {
		session = nil
		fmt.Println("too many errors, giving up submit")
	}

	// submit
	submission, err := submit(link)
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

	if err != nil {
		sess, err := createSession()
		if err != nil {
			session = nil
			fmt.Println("too many errors, giving up submit")
			return
		}
		errs++
		session = sess
		postLink(link)
	}
}

func approve(submission *geddit.Submission) error {
	v := url.Values{
		"id": {submission.FullID},
	}

	type generic struct {
		Json struct {
			Errors [][]string
		}
	}

	err := postBody("https://oauth.reddit.com/api/approve", v, &generic{})
	if err != nil {
		return err
	}

	return nil
}

func submit(link *Link) (*geddit.Submission, error) {

	// canonical nick -> flair text
	// flair := regexp.MustCompile(`[-_^]`).ReplaceAllString(link.Author, "")

	// Build form for POST request.
	v := url.Values{
		"title":       {fmt.Sprintf("<%s> %s", link.Author, link.Title)},
		"url":         {link.Url},
		"sr":          {pluginConfig.subreddit},
		"sendreplies": {strconv.FormatBool(false)},
		"resubmit":    {strconv.FormatBool(false)},
		"api_type":    {"json"},
		// "flair_text":  {flair},
		"kind": {"link"},
	}

	type submission struct {
		Json struct {
			Errors  [][]string
			Message string `json:"message"`
			Error   string `json:"error"`
			Data    geddit.Submission
		}
	}
	submit := &submission{}

	err := postBody("https://oauth.reddit.com/api/submit", v, submit)
	if err != nil {
		return nil, err
	}
	if submit.Json.Error != "" {
		return nil, errors.New("failed to submit: " + submit.Json.Message)
	}

	// TODO check s.Errors and do something useful?
	return &submit.Json.Data, nil
}

func postBody(link string, form url.Values, d interface{}) error {
	req, err := http.NewRequest("POST", link, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	// This is needed to avoid rate limits
	//req.Header.Set("User-Agent", o.UserAgent)

	// POST form provided
	req.PostForm = form

	if session.Client == nil {
		return errors.New("the OAuth Session lacks HTTP client! Use func (o OAuthSession) LoginAuth() to make one")
	}

	resp, err := session.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// The caller may want JSON decoded, or this could just be an update/delete request.
	if d != nil {
		err = json.Unmarshal(body, d)
		if err != nil {
			return err
		}
	}

	return nil
}
