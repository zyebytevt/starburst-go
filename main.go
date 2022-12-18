package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	streamdeck "github.com/magicmonkey/go-streamdeck"
	_ "github.com/magicmonkey/go-streamdeck/devices"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/zyebytevt/streaming-backend/addons/obs"
	"github.com/zyebytevt/streaming-backend/addons/twitch"
	"github.com/zyebytevt/streaming-backend/addons/vseeface"
)

var sd *streamdeck.StreamDeck

func loadConfigAndDefaults() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04"})

	// first set some default values
	viper.AddConfigPath(".")
	viper.SetDefault("buttons.images", "images/buttons") // location of button images

	viper.SetDefault("obs.host", "localhost") // OBS webhooks endpoint
	viper.SetDefault("obs.port", 4455)        // OBS webhooks endpoint

	// now read in config for any overrides
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		log.Warn().Msgf("Cannot read config file: %s \n", err)
	}
}

func main() {
	loadConfigAndDefaults()
	logrus.Info("Starting the streaming backend. ZyeByte sends her regards.")

	var err error
	sd, err = streamdeck.New()
	if err != nil {
		logrus.WithError(err).Error("Could not connect to the StreamDeck. Please check connectivity.")
		panic(err)
	}

	// sd.SetBrightness(60) // Create buttons for changing brightness

	//lib.NewButton(sd, 31, "Sanity!", nil, nil)

	if err := vseeface.Setup(sd); err != nil {
		logrus.WithError(err).Warning("Failed to initialize VSeeFace addon.")
	}

	if err := obs.Setup(sd); err != nil {
		logrus.WithError(err).Warning("Failed to initialize OBS addon.")
	}

	if err := twitch.Setup(sd); err != nil {
		logrus.WithError(err).Warning("Failed to initialize Twitch addon.")
	}

	// init MQTT
	/*mqtt_addon := addons.MqttThing{SD: sd}
	mqtt_addon.Init()
	mqtt_addon.Buttons()

	// init Screenshot
	screenshot_addon := addons.Screenshot{SD: sd}
	screenshot_addon.Init()
	screenshot_addon.Buttons()

	// init WindowManager
	windowmgmt_addon := addons.WindowMgmt{SD: sd}
	windowmgmt_addon.Init()
	windowmgmt_addon.Buttons()

	// set up soundcaster
	caster_addon := addons.Caster{SD: sd}
	caster_addon.Init()
	caster_addon.Buttons()



	// Nightbot (needs ngrok twitch if refresh has expired)
	nightbot_addon := addons.Nightbot{SD: sd}
	nightbot_addon.Init()
	nightbot_addon.Buttons()

	// Mute/Audio features
	mute_addon := addons.Mute{SD: sd, Button_id: 31}
	mute_addon.Init()
	mute_addon.Buttons()*/

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

	http.ListenAndServe(":7001", nil)
}
