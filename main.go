package main

import (
	"fmt"
	"net/http"
	"sync"

	streamdeck "github.com/magicmonkey/go-streamdeck"
	_ "github.com/magicmonkey/go-streamdeck/devices"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/zyebytevt/starburst-go/subsystems/dbus"
	"github.com/zyebytevt/starburst-go/subsystems/general"
	"github.com/zyebytevt/starburst-go/subsystems/obs"
	"github.com/zyebytevt/starburst-go/subsystems/twitch"
	"github.com/zyebytevt/starburst-go/subsystems/vseeface"
)

var sd *streamdeck.StreamDeck

func loadConfigAndDefaults() {
	// first set some default values
	viper.AddConfigPath(".")

	viper.SetDefault("obs.host", "localhost") // OBS webhooks endpoint
	viper.SetDefault("obs.port", 4455)        // OBS webhooks endpoint

	// now read in config for any overrides
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		logrus.WithError(err).Warning("Could not read config file.")
	}
}

func main() {
	loadConfigAndDefaults()
	logrus.Info("StarBurst, Go! ZyeByte sends her regards.")

	var err error
	sd, err = streamdeck.New()
	if err != nil {
		logrus.WithError(err).Error("Could not connect to the StreamDeck. Please check connectivity.")
		panic(err)
	}

	if err := general.Setup(sd); err != nil {
		logrus.WithError(err).Error("Failed to initialize general functionality.")
		panic(err)
	}

	if err := dbus.Setup(sd); err != nil {
		logrus.WithError(err).Warning("Failed to initialize D-Bus addon.")
	}

	if err := vseeface.Setup(sd); err != nil {
		logrus.WithError(err).Warning("Failed to initialize VSeeFace addon.")
	}

	if err := obs.Setup(sd); err != nil {
		logrus.WithError(err).Warning("Failed to initialize OBS addon.")
	}

	if err := twitch.Setup(sd); err != nil {
		logrus.WithError(err).Warning("Failed to initialize Twitch addon.")
	}

	// TODO: Do we need to have this webserver running at all times? Maybe just
	// for authentication is enough.
	go webserver()

	logrus.Info("Up and running!")
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}

func webserver() {
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	})

	http.HandleFunc("/auth-callback", twitch.HandleAuthCallback)

	logrus.Info("Starting up webserver...")
	http.ListenAndServe(":7001", nil)
}
