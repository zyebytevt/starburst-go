package addons

import (
	"fmt"
	"image/color"

	// "os/exec"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/events"
	"github.com/andreykaipov/goobs/api/requests/scenes"
	"github.com/magicmonkey/go-streamdeck"
	"github.com/magicmonkey/go-streamdeck/buttons"
	sddecorators "github.com/magicmonkey/go-streamdeck/decorators"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Obs struct {
	obs_client *goobs.Client

	StreamDeck *streamdeck.StreamDeck
	Offset     int
}

var active_scene_border = sddecorators.NewBorder(5, color.RGBA{255, 0, 0, 255})
var obs_scene_buttons []*ObsSceneButton

type ObsSceneButton struct {
	SceneName  string `mapstructure:"scene_name"`
	ButtonText string `mapstructure:"button_text"`
	Position   int    `mapstructure:"position"`
}

func (o *Obs) Init() {
	o.ConnectOBS()

	go o.obs_client.Listen(func(event any) {
		switch e := event.(type) {
		case *events.InputVolumeChanged:
			log.Debug().Msgf("%s's volume is now %f", e.InputName, e.InputVolumeDb)

		case *events.CurrentProgramSceneChanged:
			log.Debug().Msgf("Program scene changed event to %s.", e.SceneName)
			for _, button := range obs_scene_buttons {
				if button.SceneName == e.SceneName {
					o.StreamDeck.SetDecorator(button.Position, active_scene_border)
				} else {
					o.StreamDeck.UnsetDecorator(button.Position)
				}
			}
		}
	})
}

func (o *Obs) ConnectOBS() {
	log.Debug().Msg("Connecting to OBS...")

	var err error
	o.obs_client, err = goobs.New(fmt.Sprintf("%s:%d", viper.GetString("obs.host"), viper.GetInt("obs.port")),
		goobs.WithPassword(viper.GetString("obs.password")))

	if err != nil {
		log.Warn().Err(err).Msg("Cannot connect to OBS")
		panic(err)
	}
}

/*
func (o *Obs) ObsEventHandlers() {
	if o.obs_client.Connected() {
		// Scene change
		o.obs_client.AddEventHandler("SwitchScenes", func(e obsws.Event) {
			// Make sure to assert the actual event type.
			scene := strings.ToLower(e.(obsws.SwitchScenesEvent).SceneName)
			log.Info().Msg("Old scene: " + obs_current_scene)
			// undecorate the old
			if oldb, ok := buttons_obs[obs_current_scene]; ok {
				log.Info().Int("button", oldb.ButtonId).Msg("Clear original button decoration")
				o.SD.UnsetDecorator(oldb.ButtonId)
			}
			// decorate the new
			log.Info().Msg("New scene: " + scene)
			if eventb, ok := buttons_obs[scene]; ok {
				decorator2 := sddecorators.NewBorder(5, color.RGBA{255, 0, 0, 255})
				log.Info().Int("button", eventb.ButtonId).Msg("Highlight new scene button")
				o.SD.SetDecorator(eventb.ButtonId, decorator2)
			}
			obs_current_scene = scene
		})

		// OBS Exits
		o.obs_client.AddEventHandler("Exiting", func(e obsws.Event) {
			log.Info().Msg("OBS has exited")
			o.ClearButtons()
		})

		// Scene Collection Switched
		o.obs_client.AddEventHandler("SceneCollectionChanged", func(e obsws.Event) {
			log.Info().Msg("Scene collection changed")
			o.ClearButtons()
			o.Buttons()
		})

	}
}*/

func (o *Obs) CreateButtons() {
	// OBS Scenes to Buttons
	obs_scene_buttons = make([]*ObsSceneButton, 0)
	viper.UnmarshalKey("obs.scenes", &obs_scene_buttons)

	scenes, err := o.obs_client.Scenes.GetSceneList()
	if err != nil {
		log.Warn().Err(err)
	}

	// make buttons for these scenes
	for _, button_info := range obs_scene_buttons {
		log.Debug().Msg("Scene: " + button_info.SceneName)
		oaction := &OBSSceneAction{Scene: button_info.SceneName, Obs: o}
		/*sceneName := strings.ToLower(scene.SceneName)

		if s, ok := buttons_obs[sceneName]; ok {
			if s.Image != "" {
				image = filepath.Join(image_path, s.Image)
			}
		} else {
			// there wasn't an entry in the buttons for this scene so add one
			//buttons_obs[sceneName] = &ObsScene{}
			log.Warn().Msgf("Button for non-existent %s defined!", sceneName)
			continue
		}

		if image != "" {
			// try to make an image button

			obutton, err := buttons.NewImageFileButton(image)
			if err == nil {
				obutton.SetActionHandler(oaction)
				o.SD.AddButton(i+o.Offset, obutton)
				// store which button we just set
				buttons_obs[sceneName].SetButtonId(i + o.Offset)
			} else {
				// something went wrong with the image, use a default one
				image = image_path + "/play.jpg"
				obutton, err := buttons.NewImageFileButton(image)
				if err == nil {
					obutton.SetActionHandler(oaction)
					o.SD.AddButton(i+o.Offset, obutton)
					// store which button we just set
					buttons_obs[sceneName].SetButtonId(i + o.Offset)
				}
			}
		} else {
			// use a text button
			oopbutton := buttons.NewTextButton(scene.Name)
			oopbutton.SetActionHandler(oaction)
			o.SD.AddButton(i+o.Offset, oopbutton)
			// store which button we just set
			buttons_obs[sceneName].SetButtonId(i + o.Offset)
		}*/

		oopbutton := buttons.NewTextButton(button_info.ButtonText)
		oopbutton.SetActionHandler(oaction)
		o.StreamDeck.AddButton(button_info.Position, oopbutton)

		if button_info.SceneName == scenes.CurrentProgramSceneName {
			o.StreamDeck.SetDecorator(button_info.Position, active_scene_border)
		}
	}

	// highlight the active scene

	/*if eventb, ok := buttons_obs[scenes.CurrentProgramSceneName]; ok {
		decorator2 := sddecorators.NewBorder(5, color.RGBA{255, 0, 0, 255})
		log.Info().Int("button", eventb.ButtonId).Msg("Highlight current scene")
		o.SD.SetDecorator(eventb.ButtonId, decorator2)
	}*/

	// show a button to reinitialise all the OBS things
	//startbutton := buttons.NewTextButton("Go OBS")
	//startbutton.SetActionHandler(&OBSStartAction{Client: o.obs_client, Obs: o})
	//o.SD.AddButton(o.Offset+7, startbutton)
}

/*func (o *Obs) ClearButtons() {
	for i := 0; i < 7; i++ {
		o.StreamDeck.UnsetDecorator(o.Offset + i)
		clearbutton := buttons.NewTextButton("")
		o.StreamDeck.AddButton(o.Offset+i, clearbutton)
	}
}*/

type OBSSceneAction struct {
	Obs   *Obs
	Scene string
	btn   streamdeck.Button
}

func (action *OBSSceneAction) Pressed(btn streamdeck.Button) {
	log.Info().Msg("Set scene: " + action.Scene)

	_, err := action.Obs.obs_client.Scenes.SetCurrentProgramScene(&scenes.SetCurrentProgramSceneParams{SceneName: action.Scene})
	if err != nil {
		log.Error().Err(err).Msg("Failed to set scene.")
	}
}
