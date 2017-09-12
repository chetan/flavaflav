package btc

import (
	"fmt"

	"github.com/go-chat-bot/bot"
)

var (
	price = CoindeskSource{}
)

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
