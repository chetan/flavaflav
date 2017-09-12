package url

import (
	"fmt"
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

	fmt.Println("res:", res)
	if res != `[ https://99c.org/gEI ] Jeff Piotrowski: "Hurricane Irma eye wall. #flwx #hutticaneIrma"` {
		t.Fail()
	}

}
