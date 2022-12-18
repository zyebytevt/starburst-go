package obs

import (
	"github.com/andreykaipov/goobs/api/requests/scenes"
	"github.com/sirupsen/logrus"
	"github.com/zyebytevt/streaming-backend/lib"
)

func setSceneCallback(config *lib.ButtonConfig) {
	_, err := obsClient.Scenes.SetCurrentProgramScene(&scenes.SetCurrentProgramSceneParams{SceneName: config.Parameters["scene_name"].(string)})
	if err != nil {
		logrus.WithError(err).Error("Failed to set scene!")
	}
}
