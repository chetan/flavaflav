package instagram

import (
	"fmt"
	"testing"
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
