package lib

import (
	"github.com/magicmonkey/go-streamdeck"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ButtonConfig struct {
	ActionName  string         `mapstructure:"action"`
	Parameters  map[string]any `mapstructure:"params"`
	ButtonText  string         `mapstructure:"button_text"`
	ButtonIndex int            `mapstructure:"button_index"`
}

func GetConfigsForKey(key string) ([]*ButtonConfig, error) {
	configs := make([]*ButtonConfig, 0)
	err := viper.UnmarshalKey(key, &configs)

	return configs, err
}

func CreateButtonFromConfig(streamDeck *streamdeck.StreamDeck, config *ButtonConfig, callbacks map[string]ActionCallback) *Button {
	callback, exists := callbacks[config.ActionName]
	if !exists {
		logrus.Warningf("Action callback '%s' does not exist.", config.ActionName)
	}

	return NewButton(streamDeck, config.ButtonIndex, config.ButtonText, callback, config)
}
