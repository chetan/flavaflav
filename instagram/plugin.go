package instagram

import (
	"fmt"
	"strings"
	"time"

	"github.com/Vorkytaka/instagram-go-scraper/instagram"
	"github.com/chetan/flavaflav/util"
	"github.com/go-chat-bot/bot"
)

func init() {
	bot.RegisterPassiveCommand(
		"instagram",
		handleInstagram)
}

func fetchInsta(u string) (*instagram.Media, error) {
	media, err := instagram.GetMediaByURL(u)
	if err != nil {
		return nil, err
	}
	return &media, nil
}

func handleInstagram(cmd *bot.PassiveCmd) (string, error) {

	if util.IgnoreCmd(cmd) {
		return "", nil
	}

	URL := util.ExtractURL(cmd.Raw)
	if URL == "" {
		return "", nil
	}

	if util.InstagramRe.MatchString(URL) {
		insta, err := fetchInsta(URL)
		if err != nil {
			return "", err
		}

		t := time.Unix(int64(insta.Date), 0)
		caption := insta.Caption
		lines := strings.Split(caption, "\n")
		if len(lines) > 1 {
			caption = lines[0] + " ..."
		}
		return fmt.Sprintf("<%s> %s // %s", insta.Owner.Username, caption, t.String()), nil
	}

	return "", nil
}
