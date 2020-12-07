package covid

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-chat-bot/bot"
	colly "github.com/gocolly/colly/v2"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	currentStatsURL = "https://api.covidtracking.com/v1/us/current.json"
)

// Enable the covid command plugin
func Enable() {
	bot.RegisterCommand("covid", "Get current covid stats for US", "", getCovidStats)
}

func fetchCovidStatus(results chan string) error {
	c := colly.NewCollector(colly.AllowedDomains("www.worldometers.info"))
	c.OnHTML("table#usa_table_countries_today", func(table *colly.HTMLElement) {
		tr := table.DOM.Find("tbody tr:first-of-type")
		tds := tr.Children()
		newCases := tds.Eq(3).Text()
		newDeaths := tds.Eq(5).Text()
		results <- fmt.Sprintf("today: new cases = %s; new deaths = %s", newCases, newDeaths)
	})
	c.OnHTML("table#usa_table_countries_yesterday", func(table *colly.HTMLElement) {
		tr := table.DOM.Find("tbody tr:first-of-type")
		tds := tr.Children()
		newCases := tds.Eq(3).Text()
		newDeaths := tds.Eq(5).Text()
		results <- fmt.Sprintf("yesterday: new cases = %s; new deaths = %s", newCases, newDeaths)
	})
	return c.Visit("https://www.worldometers.info/coronavirus/country/us/")
}

func getCovidStats(cmd *bot.Cmd) (string, error) {
	resultsChan := make(chan string, 1)
	go fetchCovidStatus(resultsChan)

	out := <-resultsChan
	out += "\n" + <-resultsChan
	return out, nil

}

func fetchStatusApi() (string, error) {
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
