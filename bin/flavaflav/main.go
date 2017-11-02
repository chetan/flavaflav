package main

import (
	"fmt"
	"os/signal"
	"strings"
	"syscall"

	"github.com/chetan/flavaflav/trumpykins"

	"github.com/spf13/viper"

	"github.com/go-chat-bot/bot/irc"
	_ "github.com/go-chat-bot/plugins/catgif"
	_ "github.com/go-chat-bot/plugins/chucknorris"
	// _ "github.com/go-chat-bot/plugins/url"

	"github.com/chetan/flavaflav/btc"
	_ "github.com/chetan/flavaflav/instagram"
	_ "github.com/chetan/flavaflav/twitter"
	_ "github.com/chetan/flavaflav/url"
	"github.com/chetan/flavaflav/util"

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

	util.IgnoreNicks = viper.GetStringSlice("ignore_nicks")
	util.IgnorePatterns = viper.GetStringSlice("ignore_patterns")

	// btc plugin
	btcChannels := viper.GetStringSlice("btc_channels")
	if len(btcChannels) > 0 {
		fmt.Println("enabling btc monitor")
		btc.SetChannels(btcChannels)
	}

	// trumpykins plugin
	twitterKey := viper.GetString("twitter_key")
	twitterSecret := viper.GetString("twitter_secret")
	twitterAccessToken := viper.GetString("twitter_access_token")
	twitterAccessSecret := viper.GetString("twitter_access_secret")
	trumpChannels := viper.GetStringSlice("trump_channels")
	if twitterKey != "" && twitterSecret != "" && len(trumpChannels) > 0 {
		fmt.Println("enabling trumpykins stream")
		trumpykins.Enable(twitterKey, twitterSecret, twitterAccessToken, twitterAccessSecret, trumpChannels)
	}

	server := viper.GetString("server")
	if !strings.Contains(server, ":") {
		server += ":6667" // append default port
	}

	cfg := &irc.Config{
		Server:   server,
		Channels: viper.GetStringSlice("channels"),
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

	trap()

	irc.Run(cfg)

}

func trap() {
	sigs := make(chan os.Signal, 1)
	go func() {
		s := <-sigs // blocks until signal received
		fmt.Println("Caught signal: ", s)
		os.Exit(0)
	}()
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGQUIT)
}
