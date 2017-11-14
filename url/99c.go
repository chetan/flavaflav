package url

import (
	"net/url"
	"strings"

	"github.com/chetan/flavaflav/util"
)

type Response struct {
	ShortURL string `json:"shorturl"`
}

func shortenURL99c(u string) (string, error) {
	res := Response{}

	encodedURL := url.QueryEscape(u)
	err := util.GetJSON("https://99c.org/botapi.php?action=shorturl&format=json&url="+encodedURL, &res)
	if err != nil {
		return "", err
	}

	return strings.Replace(res.ShortURL, "http://", "https://", 1), nil
}
