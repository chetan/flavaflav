package cloneit

import (
	"fmt"
	"sync"

	"github.com/jzelinskie/geddit"
	"github.com/pkg/errors"
)

type config struct {
	username  string
	password  string
	subreddit string
	channel   string
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

var session *geddit.LoginSession

func AddLink(link *Link) {
	if enabled && link.Channel == pluginConfig.channel {
		go postLink(link) // submit in background
	}
}

// Enable the plugin
func Enable(user string, pass string, sub string, channel string) {
	pluginConfig = &config{user, pass, sub, channel}
	sess, err := createSession()
	if err != nil {
		fmt.Println("cloneit: err creating initial session: ", err)
		return
	}
	session = sess
	enabled = true
}

func createSession() (*geddit.LoginSession, error) {
	if errs > 2 {
		// give up
		return nil, errors.New("too many errors, giving up creating session")
	}

	sess, err := geddit.NewLoginSession(
		pluginConfig.username,
		pluginConfig.password,
		"flavaflav v1",
	)

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

	err := session.Submit(&geddit.NewSubmission{
		Subreddit: pluginConfig.subreddit,
		Title:     fmt.Sprintf("<%s> %s", link.Author, link.Title),
		Content:   link.Url,
		Captcha:   &geddit.Captcha{},
	})

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
