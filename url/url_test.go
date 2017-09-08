package url

import "testing"

func TestShortenURL(t *testing.T) {
	s, err := shortenURL("http://yahoo.com")
	if err != nil {
		t.Fail()
	}
	if s != "http://tinyurl.com/2ks" {
		t.Fail()
	}
}
