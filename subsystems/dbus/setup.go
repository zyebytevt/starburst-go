package dbus

import (
	"github.com/godbus/dbus"
	"github.com/magicmonkey/go-streamdeck"
	"github.com/sirupsen/logrus"
	"github.com/zyebytevt/starburst-go/lib"
)

var connection *dbus.Conn

var actionCallbacks map[string]lib.ActionCallback = map[string]lib.ActionCallback{
	"call": callCallback,
}

func Setup(streamDeck *streamdeck.StreamDeck) error {
	logrus.Info("Initializing D-Bus addon...")

	var err error

	connection, err = dbus.SessionBus()

	if err != nil {
		return err
	}

	configs, err := lib.GetConfigsForKey("dbus.buttons")
	if err != nil {
		return err
	}

	for _, config := range configs {
		lib.CreateButtonFromConfig(streamDeck, config, actionCallbacks)
	}

	return err
}
