package twitter

import (
	"fmt"
	"testing"
)

func TestShortenURL(t *testing.T) {

	// z := "abc aaa def aaa hij"
	// i := strings.LastIndex(z, "aaa")
	// q := z[0:i] + z[i+3:]
	// println("'", q, "'")
	//
	// t.FailNow()

	tweet, err := fetchTweet("https://mobile.twitter.com/realDonaldTrump/status/905958330815926276")
	if err != nil {
		t.Fail()
	}
	fmt.Println(tweet)
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
