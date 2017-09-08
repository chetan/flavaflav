package url

import (
	"net/url"

	"github.com/go-chat-bot/plugins/web"
)

func shortenURLTiny(u string) (string, error) {
	encodedURL := url.QueryEscape(u)
	body, err := web.GetBody("http://tinyurl.com/api-create.php?url=" + encodedURL)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
