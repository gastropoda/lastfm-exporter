package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shkh/lastfm-go/lastfm"
	"github.com/spf13/viper"
)

var lastfmTrackPlays = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "lastfm_track_plays",
	Help: "Count each play of a track",
}, []string{
	"artist",
	"album",
	"track",
})

func syncStats(api *lastfm.Api, user string) {
	from := 0
	go func() {
		for {
			result, _ := api.User.GetRecentTracks(lastfm.P{
				"user": user,
				"from": from + 1,
			})
			for _, track := range result.Tracks {
				if track.NowPlaying != "true" {
					continue
				}

				lastfmTrackPlays.With(prometheus.Labels{
					"artist": track.Artist.Name,
					"album":  track.Album.Name,
					"track":  track.Name,
				}).Set(1)
			}
			time.Sleep(5 * time.Second)
			lastfmTrackPlays.Reset()
		}
	}()
}

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.config/lastfm-exporter")
	viper.AddConfigPath(".")

	viper.SetDefault("ApiKey", "a60a760a91999658f00cc128f4b17100")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	api_key := viper.GetString("lastfm.api.key")
	api_secret := viper.GetString("lastfm.api.secret")
	user := viper.GetString("lastfm.user")

	api := lastfm.New(api_key, api_secret)
	syncStats(api, user)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
