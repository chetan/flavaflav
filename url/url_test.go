package url

import (
	"testing"

	"github.com/go-chat-bot/bot"
	"github.com/stretchr/testify/assert"
)

func TestShortenURL(t *testing.T) {
	s, err := shortenURL("http://yahoo.com")
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
	title, err := extractTitle(u)
	assert.NoError(t, err)
	assert.Equal(t, `The Canadian Ghost Town That Tesla Is Bringing Back to Life - Bloomberg`, title)

}

func TestJalopnikTitle(t *testing.T) {
	u := "https://blackflag.jalopnik.com/what-108-years-of-repaving-looks-like-under-indianapoli-1820048121"

	cmd := bot.PassiveCmd{Raw: u}
	res, err := urlTitle(&cmd)
	assert.NoError(t, err)
	assert.Equal(t, `[ https://99c.org/gL5 ] What 108 Years Of Repaving Looks Like Under Indianapolis Motor Speedway's Asphalt`, res)
}
