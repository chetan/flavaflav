package covid

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-chat-bot/bot"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	currentStatsURL = "https://api.covidtracking.com/v1/us/current.json"
)

// Enable the covid command plugin
func Enable() {
	fmt.Println("adding covid command")
	bot.RegisterCommand("covid", "Get current covid stats for US", "", getCovidStats)
}

func getCovidStats(cmd *bot.Cmd) (string, error) {

	b, err := fetch(currentStatsURL)
	if err != nil {
		return "", err
	}
	var data []map[string]interface{}
	json.Unmarshal(b, &data)

	p := message.NewPrinter(language.English)
	d := data[0]
	out := p.Sprintf("new positive tests: %.0f // new deaths: %.0f", d["positiveIncrease"], d["deathIncrease"])

	yest := time.Now().Add(-time.Hour * 24).Format("20060102")
	if yest == fmt.Sprintf("%.0f", d["date"]) {
		// still return yesterday's data, just go with that
		return "yesterday: " + out, nil
	}

	url := "https://api.covidtracking.com/v1/us/" + yest + ".json"
	fmt.Println("fetching yest", url)

	b, err = fetch(url)
	if err != nil {
		return "", err
	}
	var yestData map[string]interface{}
	json.Unmarshal(b, &yestData)
	out = p.Sprintf("today: %s\nyesterday: new positive tests: %.0f // new deaths: %.0f", out, d["positiveIncrease"], d["deathIncrease"])

	return out, nil
}

func fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}
