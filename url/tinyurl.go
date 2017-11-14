package url

import (
	"net/url"

	"github.com/chetan/flavaflav/util"
)

func shortenURLTiny(u string) (string, error) {
	encodedURL := url.QueryEscape(u)
	body, err := util.GetBody("http://tinyurl.com/api-create.php?url=" + encodedURL)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
