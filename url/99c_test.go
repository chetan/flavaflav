package url

import (
	"fmt"
	"testing"
)

func TestShortenURL99c(t *testing.T) {
	s, err := shortenURL99c("http://yahoo.com")
	if err != nil {
		fmt.Println("Error running 99c shortener: ", err)
		t.Fail()
	}
	if s != "https://99c.org/gEd" {
		t.Fail()
	}
}
