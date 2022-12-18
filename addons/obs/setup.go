package obs

import (
	"fmt"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/events"
	"github.com/andreykaipov/goobs/api/requests/scenes"
	"github.com/magicmonkey/go-streamdeck"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/zyebytevt/streaming-backend/lib"
)

type SceneConfig struct {
	SceneName   string `mapstructure:"scene_name"`
	ButtonText  string `mapstructure:"button_text"`
	ButtonIndex int    `mapstructure:"button_index"`
}

var obsClient *goobs.Client
var sceneButtons map[string]*lib.Button = make(map[string]*lib.Button)

func Setup(streamDeck *streamdeck.StreamDeck) error {
	logrus.Info("Initializing OBS addon...")

	var err error

	obsClient, err = goobs.New(fmt.Sprintf("%s:%d", viper.GetString("obs.host"), viper.GetInt("obs.port")),
		goobs.WithPassword(viper.GetString("obs.password")))

	if err != nil {
		return err
	}

	configs := make([]*SceneConfig, 0)
	viper.UnmarshalKey("obs.scenes", &configs)

	for _, config := range configs {
		btn := lib.NewButton(streamDeck, config.ButtonIndex, config.ButtonText, func(userData any) {
			funcConfig := userData.(*SceneConfig)

			_, err := obsClient.Scenes.SetCurrentProgramScene(&scenes.SetCurrentProgramSceneParams{SceneName: funcConfig.SceneName})
			if err != nil {
				logrus.WithError(err).Error("Failed to set scene!")
			}
		}, config)

		sceneButtons[config.SceneName] = btn
	}

	currentScene, err := obsClient.Scenes.GetCurrentProgramScene()
	if button, exists := sceneButtons[currentScene.CurrentProgramSceneName]; exists {
		button.SetActive(true)
	}

	go obsClient.Listen(obsEventListener)

	return err
}

func obsEventListener(event any) {
	switch e := event.(type) {
	case *events.InputVolumeChanged:
		logrus.Debugf("%s's volume is now %f", e.InputName, e.InputVolumeDb)

	case *events.CurrentProgramSceneChanged:
		logrus.Debugf("Program scene changed event to %s.", e.SceneName)

		for _, button := range sceneButtons {
			button.SetActive(false)
		}

		if button, exists := sceneButtons[e.SceneName]; exists {
			button.SetActive(true)
		}
	}
}
