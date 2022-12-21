package obs

import (
	"fmt"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/events"
	"github.com/magicmonkey/go-streamdeck"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/zyebytevt/starburst-go/lib"
)

var obsClient *goobs.Client
var sceneButtons map[string]*lib.Button = make(map[string]*lib.Button)
var sourceButtons map[float64]*lib.Button = make(map[float64]*lib.Button)

var actionCallbacks map[string]lib.ActionCallback = map[string]lib.ActionCallback{
	"set_scene":                setSceneCallback,
	"toggle_source_visibility": toggleSourceVisibilityCallback,
}

func Setup(streamDeck *streamdeck.StreamDeck) error {
	logrus.Info("Initializing OBS addon...")

	var err error

	obsClient, err = goobs.New(fmt.Sprintf("%s:%d", viper.GetString("obs.host"), viper.GetInt("obs.port")),
		goobs.WithPassword(viper.GetString("obs.password")))

	if err != nil {
		return err
	}

	configs, err := lib.GetConfigsForKey("obs.buttons")
	if err != nil {
		return err
	}

	for _, config := range configs {
		btn, _ := lib.CreateButtonFromConfig(streamDeck, config, actionCallbacks)

		if config.ActionName == "set_scene" {
			sceneButtons[config.Parameters["scene_name"].(string)] = btn
		} else if config.ActionName == "toggle_source_visibility" {
			id, err := getSceneItemId(config.Parameters["scene_name"].(string), config.Parameters["source_name"].(string))

			if err == nil {
				sourceButtons[id] = btn
				btn.UserData["scene_item_id"] = id

				visible, _ := getSceneItemVisibility(config.Parameters["scene_name"].(string), id)

				if visible {
					btn.SetDecorator(lib.ActiveStateDecorator)
				}
			} else {
				logrus.WithError(err).Error("Failed to set up source button.")
			}
		}
	}

	currentScene, err := obsClient.Scenes.GetCurrentProgramScene()
	if button, exists := sceneButtons[currentScene.CurrentProgramSceneName]; exists {
		button.SetDecorator(lib.ActiveStateDecorator)
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
			button.SetDecorator(nil)
		}

		// TODO: This sometimes causes a nil-pointer access (Segfault), investigate.
		if button, exists := sceneButtons[e.SceneName]; exists {
			button.SetDecorator(lib.ActiveStateDecorator)
			button.FlashNotify(lib.NotifySuccess)
		}

	case *events.SceneItemEnableStateChanged:
		if button, exists := sourceButtons[e.SceneItemId]; exists {
			if e.SceneItemEnabled {
				button.SetDecorator(lib.ActiveStateDecorator)
			} else {
				button.SetDecorator(nil)
			}
		}
	}
}
