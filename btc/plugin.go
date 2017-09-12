package btc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chat-bot/bot"
)

var (
	price = CoindeskSource{}
)

type CoindeskSource struct {
	Bpi struct {
		USD struct {
			RateFloat float64 `json:"rate_float"`
		} `json:"USD"`
		EUR struct {
			RateFloat float64 `json:"rate_float"`
		} `json:"EUR"`
	} `json:"bpi"`
}

func (s *CoindeskSource) GetUSD() float64 {
	return s.Bpi.USD.RateFloat
}

func (s *CoindeskSource) GetEUR() float64 {
	return s.Bpi.EUR.RateFloat
}

func (s *CoindeskSource) Update() error {

	client := http.Client{
		Timeout: 60 * time.Second,
	}

	resp, err := client.Get("https://api.coindesk.com/v1/bpi/currentprice.json")
	if err != nil {
		return fmt.Errorf("Failed to fetch from coindesk: %v", err)
	}
	defer resp.Body.Close()

	src := CoindeskSource{}
	err = json.NewDecoder(resp.Body).Decode(&src)
	if err != nil {
		return fmt.Errorf("Unexpected Response: %v", err)
	}

	s.Bpi = src.Bpi

	return nil
}

// SetChannels enables the BTC cron for the given list of channels
func SetChannels(channels []string) {
	cron := bot.PeriodicConfig{
		CronSpec: "@every 6h",
		Channels: channels,
		CmdFunc:  handleBTC,
	}
	bot.RegisterPeriodicCommand("btc", cron)
}

func handleBTC(channel string) (string, error) {
	err := price.Update()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("BTC $%.2f", price.GetUSD()), nil
}
