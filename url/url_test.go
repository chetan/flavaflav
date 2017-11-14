package url

import (
	"testing"

	"github.com/go-chat-bot/bot"
	"github.com/stretchr/testify/assert"
)

func TestShortenURL(t *testing.T) {
	s, err := ShortenURL("http://yahoo.com")
	if err != nil {
		t.Fail()
	}
	if s != "https://99c.org/gEd" {
		t.Fail()
	}
}

func TestFetchTitleExtraAttrs(t *testing.T) {
	// Title tag has attributes
	u := "https://www.pscp.tv/w/bIQ_ZDIzNTAxNTR8MUx5eEJFeWRsRG5KTtA5noFIAWfUnvRhOXuAiH3QgDfv-GkHkcjtUpWEfnE6"

	cmd := bot.PassiveCmd{Raw: u}
	res, err := urlTitle(&cmd)
	assert.NoError(t, err)
	assert.Equal(t, `[ https://99c.org/gEI ] Jeff Piotrowski: "Hurricane Irma eye wall. #flwx #hutticaneIrma"`, res)
}

func TestSkipTitleForLargeFiles(t *testing.T) {
	u := "http://mirror.math.princeton.edu/pub/ubuntu-iso/17.10/ubuntu-17.10-desktop-amd64.iso"

	cmd := bot.PassiveCmd{Raw: u}
	res, err := urlTitle(&cmd)
	assert.NoError(t, err)
	assert.Equal(t, `[ https://99c.org/gL6 ] application/octet-stream; Content-Length: 1.4G`, res)
}

func TestExtractTitle(t *testing.T) {
	u := "https://www.bloomberg.com/news/features/2017-10-31/the-canadian-ghost-town-that-tesla-is-bringing-back-to-life?cmpid=socialflow-twitter-business"
	ti := `[ https://99c.org/gLh ] The Canadian Ghost Town That Tesla Is Bringing Back to Life - Bloomberg`
	testUrl(t, u, ti)
}

func TestJalopnikTitle(t *testing.T) {
	u := "https://blackflag.jalopnik.com/what-108-years-of-repaving-looks-like-under-indianapoli-1820048121"
	ti := `[ https://99c.org/gL5 ] What 108 Years Of Repaving Looks Like Under Indianapolis Motor Speedway's Asphalt`
	testUrl(t, u, ti)
}

func TestAPTitle(t *testing.T) {
	u := "https://apnews.com/9a605019eeba4ad2934741091105de42"
	ti := `[ https://99c.org/gLg ] APNewsBreak: Gov't won't pursue talking car mandate`
	testUrl(t, u, ti)
}

func TestVoxTitle(t *testing.T) {
	u := "https://www.vox.com/science-and-health/2017/11/2/16594408/great-pyramid-giza-cosmic-rays-void-particle-physics-nature?utm_campaign=vox&utm_content=chorus&utm_medium=social&utm_source=twitter"
	ti := `[ https://99c.org/gLc ] Great Pyramid: Scientists found a mysterious void inside using cosmic rays - Vox`
	testUrl(t, u, ti)
}

func TestRedditTitle(t *testing.T) {
	u := "https://www.reddit.com/r/aws/comments/7cufe4/amazon_web_services_denies_reports_of_china_exit/"
	ti := `[ https://99c.org/gOZ ] Amazon Web Services denies reports of China exit, confirms some asset sales : aws`
	testUrl(t, u, ti)
}

func testUrl(t *testing.T, url string, expectedTitle string) {
	cmd := bot.PassiveCmd{Raw: url}
	title, err := urlTitle(&cmd)
	assert.NoError(t, err)
	assert.Equal(t, expectedTitle, title)
}
