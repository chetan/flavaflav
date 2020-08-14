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

type config struct {
	consumerKey    string
	consumerSecret string
	accessToken    string
	accessSecret   string
	channels       []string
}

var (
	pendingTweets = []string{}
	mtx           = sync.Mutex{}
	pluginConfig  = &config{}
)

var twitterIDs = []int64{
	25073877, // realDonaldTrump
	939091,   // JoeBiden
	30354991, // KamalaHarris
}

var twitterIDmap map[int64]int

func Enable(consumerKey string, consumerSecret string, accessToken string, accessSecret string, channels []string) {
	pluginConfig = &config{consumerKey, consumerSecret, accessToken, accessSecret, channels}
	startPlugin()

	cron := bot.PeriodicConfig{
		CronSpec: "@every 1m",
		Channels: channels,
		CmdFunc:  postTrumpTweets,
	}
	bot.RegisterPeriodicCommand("trumpykins", cron)
}

func startPlugin() {
	fmt.Println("starting trumpykins tracker")

	// client setup
	config := oauth1.NewConfig(pluginConfig.consumerKey, pluginConfig.consumerSecret)
	token := oauth1.NewToken(pluginConfig.accessToken, pluginConfig.accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	// create ID params for listening
	twitterIDmap = make(map[int64]int)
	var followIDs []string
	for _, id := range twitterIDs {
		twitterIDmap[id] = 1
		followIDs = append(followIDs, fmt.Sprintf("%d", id))
	}

	// create stream
	params := &twitter.StreamFilterParams{
		Follow:        followIDs,
		StallWarnings: twitter.Bool(true),
	}
	stream, err := client.Streams.Filter(params)
	if err != nil {
		fmt.Println("failed to setup twitter streaming api: ", err)
		return
	}

	// setup listener
	demux := twitter.NewSwitchDemux()
	demux.Tweet = handleTweet

	// start runloop
	go func() {
		defer stream.Stop()
		demux.HandleChan(stream.Messages) // runs forever
		fmt.Println("uh oh, demux exited the loop... restarting")
		startPlugin()
	}()
}

/**
 * Handle incoming tweet from firehose stream created in startPlugin.
 */
func handleTweet(tweet *twitter.Tweet) {
	if !isFollowing(tweet.User.ID) {
		return // ignore
	}

	fmt.Printf("new tweet: %#v\n", tweet)

	tweetURL := fmt.Sprintf("https://twitter.com/%s/status/%s", tweet.User.ScreenName, tweet.IDStr)
	fullTweet, err := twitter_plugin.FetchTweet(tweetURL)
	if err != nil {
		fmt.Println("error fetching full tweet: ", err)
		return
	}
	fullTweet.AuthorID = tweet.User.ID
	twit := fullTweet.String()

	// add RT prefix if retweet
	if tweet.RetweetedStatus != nil {
		twit = fmt.Sprintf("<@%s> RT %s", tweet.User.ScreenName, twit)
	}

	// append short url
	shortURL, err := url.ShortenURL(tweetURL)
	if err == nil {
		twit += " // " + shortURL
	} else {
		twit += " // " + tweetURL // just use long one on err
	}

	// fmt.Println(twit)

	// add new tweet to our buffer
	mtx.Lock()
	defer mtx.Unlock()
	pendingTweets = append(pendingTweets, twit)
}

func isFollowing(id int64) bool {
	return twitterIDmap[id] == 1
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
