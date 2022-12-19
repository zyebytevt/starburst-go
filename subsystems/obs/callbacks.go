package obs

import (
	"github.com/andreykaipov/goobs/api/requests/scenes"
	"github.com/sirupsen/logrus"
	"github.com/zyebytevt/starburst-go/lib"
)

func setSceneCallback(button *lib.Button) {
	button.SetHighlight(lib.HighlightInProgress)

	_, err := obsClient.Scenes.SetCurrentProgramScene(&scenes.SetCurrentProgramSceneParams{SceneName: button.Config.Parameters["scene_name"].(string)})
	if err != nil {
		logrus.WithError(err).Error("Failed to set scene!")
	}
}
