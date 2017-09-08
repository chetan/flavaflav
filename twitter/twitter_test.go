package twitter

import (
	"encoding/json"
	"fmt"
	"testing"

	yaml "gopkg.in/yaml.v1"
)

func TestFetchTweet(t *testing.T) {

	tweet, err := fetchTweet("https://mobile.twitter.com/realDonaldTrump/status/905958330815926276")
	if err != nil {
		t.Fail()
	}
	// fmt.Println(tweet)

	if tweet == nil {
		t.Fail()
	}
	if tweet.AuthorName != "Donald J. Trump" {
		t.Fail()
	}

	if tweet.Text == "" {
		t.Fail()
	}

}

func TestProcessing(t *testing.T) {

	input := `
{"url":"https:\/\/twitter.com\/realDonaldTrump\/status\/905958330815926276","author_name":"Donald J. Trump","author_url":"https:\/\/twitter.com\/realDonaldTrump","html":"\u003Cblockquote class=\"twitter-tweet\"\u003E\u003Cp lang=\"en\" dir=\"ltr\"\u003EI encourage EVERYONE in the path of \u003Ca href=\"https:\/\/twitter.com\/hashtag\/HurricaneIrma?src=hash\"\u003E#HurricaneIrma\u003C\/a\u003E to heed the advice and orders of local &amp; state officials! \u003Ca href=\"https:\/\/t.co\/AQmawTpZs0\"\u003Ehttps:\/\/t.co\/AQmawTpZs0\u003C\/a\u003E\u003C\/p\u003E&mdash; Donald J. Trump (@realDonaldTrump) \u003Ca href=\"https:\/\/twitter.com\/realDonaldTrump\/status\/905958330815926276\"\u003ESeptember 8, 2017\u003C\/a\u003E\u003C\/blockquote\u003E\n\u003Cscript async src=\"\/\/platform.twitter.com\/widgets.js\" charset=\"utf-8\"\u003E\u003C\/script\u003E","width":550,"height":null,"type":"rich","cache_age":"3153600000","provider_name":"Twitter","provider_url":"https:\/\/twitter.com","version":"1.0"}
	`

	tweet := Tweet{}
	err := json.Unmarshal([]byte(input), &tweet)
	if err != nil {
		t.Fail()
	}

	processTweet(&tweet)

	b, _ := yaml.Marshal(&tweet)
	fmt.Printf("%s\n", string(b))

	println(fmt.Sprintf("<%s> %s // %s", tweet.AuthorHandle, tweet.TextBody, tweet.TextDate))

}
