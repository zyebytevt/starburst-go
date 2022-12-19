package general

import (
	"github.com/magicmonkey/go-streamdeck"
	"github.com/sirupsen/logrus"
	"github.com/zyebytevt/starburst-go/lib"
)

var sd *streamdeck.StreamDeck
var brightness int

var actionCallbacks map[string]lib.ActionCallback = map[string]lib.ActionCallback{
	"set_brightness": setBrightnessCallback,
	"execute":        executeCallback,
}

func Setup(streamDeck *streamdeck.StreamDeck) error {
	logrus.Info("Initializing general functionality...")

	sd = streamDeck

	brightness = 50
	sd.SetBrightness(brightness)

	configs, err := lib.GetConfigsForKey("general.buttons")
	if err != nil {
		return err
	}

	for _, config := range configs {
		lib.CreateButtonFromConfig(streamDeck, config, actionCallbacks)
	}

	return err
}
