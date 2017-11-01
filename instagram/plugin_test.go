package instagram

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-chat-bot/bot"
	"github.com/stretchr/testify/assert"
)

func TestFetchInstagram(t *testing.T) {
	u := "https://www.instagram.com/p/BVkpNCiDIn4/?taken-by=rebeccablikes"
	media, err := fetchInsta(u)
	if err != nil {
		fmt.Println("error fetching insta:", err)
		t.Fail()
	}

	fmt.Println("media:", media)
}

func TestHandleInstagram(t *testing.T) {

	u := "blahblah https://www.instagram.com/p/BVkpNCiDIn4/?taken-by=rebeccablikes"
	cmd := bot.PassiveCmd{Raw: u}
	res, err := handleInstagram(&cmd)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(strings.Split(res, "\n")), "Max of 1 lines returned")
	assert.True(t, strings.Contains(res, "..."))
}
