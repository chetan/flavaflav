package main

import (
	"fmt"
	"strings"

	"github.com/go-chat-bot/bot/irc"
	_ "github.com/go-chat-bot/plugins/catgif"
	_ "github.com/go-chat-bot/plugins/chucknorris"
	// _ "github.com/go-chat-bot/plugins/url"
	"github.com/spf13/viper"

	_ "github.com/chetan/flavaflav/twitter"
	_ "github.com/chetan/flavaflav/url"

	"os"
)

func initConfig() {
	viper.SetConfigName(".flavaflav")

	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("IRC")

	err := viper.ReadInConfig() // read in config file, ignore any errors
	if err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func main() {

	initConfig()

	server := viper.GetString("server")
	if !strings.Contains(server, ":") {
		server += ":6667" // append default port
	}

	var channels []string
	channels = viper.GetStringSlice("channels")
	// channels = strings.Split(os.Getenv("IRC_CHANNELS"), ",")

	cfg := &irc.Config{
		Server:   server,
		Channels: channels,
		User:     viper.GetString("nick"), // yes, these are backwards!
		Nick:     viper.GetString("user"),
		Password: viper.GetString("password"),
		UseTLS:   viper.GetBool("tls"),
		Debug:    os.Getenv("DEBUG") != "",
	}

	if cfg.Server == "" || cfg.Channels == nil || cfg.Nick == "" {
		fmt.Println("Config not found! bye")
		os.Exit(1)
	}

	fmt.Println("Running with config:", cfg)

	irc.Run(cfg)

}
