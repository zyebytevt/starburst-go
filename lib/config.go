package lib

import (
	"fmt"

	"github.com/magicmonkey/go-streamdeck"
	"github.com/spf13/viper"
)

type ButtonConfig struct {
	ActionName  string         `mapstructure:"action"`
	Parameters  map[string]any `mapstructure:"params"`
	ButtonImage string         `mapstructure:"button_image"`
	ButtonIndex int            `mapstructure:"button_index"`
}

func GetConfigsForKey(key string) ([]*ButtonConfig, error) {
	configs := make([]*ButtonConfig, 0)
	err := viper.UnmarshalKey(key, &configs)

	return configs, err
}

func CreateButtonFromConfig(streamDeck *streamdeck.StreamDeck, config *ButtonConfig, callbacks map[string]ActionCallback) (*Button, error) {
	callback, exists := callbacks[config.ActionName]
	if !exists {
		return nil, fmt.Errorf("action callback '%s' does not exist", config.ActionName)
	}

	return NewButton(streamDeck, config.ButtonIndex, config.ButtonImage, callback, config)
}
