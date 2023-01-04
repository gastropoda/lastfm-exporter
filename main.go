package main

import (
	"fmt"

	"github.com/shkh/lastfm-go/lastfm"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.config/lastfm-exporter")
	viper.AddConfigPath(".")

	viper.SetDefault("ApiKey", "a60a760a91999658f00cc128f4b17100")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	api_key := viper.GetString("ApiKey")
	api_secret := viper.GetString("ApiSecret")

	api := lastfm.New(api_key, api_secret)
	result, _ := api.Artist.GetTopTracks(lastfm.P{"artist": "Avicii"})
	for _, track := range result.Tracks {
		fmt.Println(track.Name)
	}
}
