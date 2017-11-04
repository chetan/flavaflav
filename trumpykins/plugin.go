package trumpykins

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/go-chat-bot/bot"

	twitter_plugin "github.com/chetan/flavaflav/twitter"
	"github.com/chetan/flavaflav/url"
)

var (
	pendingTweets = []string{}
	mtx           = sync.Mutex{}
)

const twitterID = 25073877 // realDonaldTrump

func Enable(consumerKey string, consumerSecret string, accessToken string, accessSecret string, channels []string) {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	params := &twitter.StreamFilterParams{
		Follow:        []string{fmt.Sprintf("%d", twitterID)},
		StallWarnings: twitter.Bool(true),
	}
	stream, err := client.Streams.Filter(params)
	if err != nil {
		fmt.Println("failed to setup twitter streaming api: ", err)
		return
	}

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		if tweet.User.ID != twitterID {
			return // ignore
		}
		tweetURL := fmt.Sprintf("https://twitter.com/%s/status/%s", tweet.User.ScreenName, tweet.IDStr)
		fmt.Println(tweetURL)
		// twit := formatTweet(tweet)
		fullTweet, err := twitter_plugin.FetchTweet(tweetURL)
		if err != nil {
			fmt.Println("error fetching full tweet: ", err)
			return
		}
		twit := fullTweet.String()

		// append short url
		shortURL, err := url.ShortenURL(tweetURL)
		if err == nil {
			twit += " // " + shortURL
		} else {
			twit += " // " + tweetURL // just use long one on err
		}

		fmt.Println(twit)

		// add new tweet to our buffer
		mtx.Lock()
		defer mtx.Unlock()
		pendingTweets = append(pendingTweets, twit)
	}

	go func() {
		demux.HandleChan(stream.Messages) // runs forever
		fmt.Println("uh oh, demux exited the loop...")
	}()

	cron := bot.PeriodicConfig{
		CronSpec: "@every 1m",
		Channels: channels,
		CmdFunc:  postTrumpTweets,
	}
	bot.RegisterPeriodicCommand("trumpykins", cron)
}

func formatTweet(tweet *twitter.Tweet) string {
	ts, _ := tweet.CreatedAtTime()
	// var text string
	return fmt.Sprintf("<@%s> %s // %s", tweet.User.ScreenName, tweet.Text, ts.Format(time.UnixDate))
}

func postTrumpTweets(channel string) (string, error) {
	mtx.Lock()
	defer mtx.Unlock()
	if len(pendingTweets) > 0 {
		out := strings.Join(pendingTweets, "\n")
		pendingTweets = []string{} // clear the buffer
		return out, nil
	}
	return "", nil
}
