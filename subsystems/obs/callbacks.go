package obs

import (
	"github.com/andreykaipov/goobs/api/requests/scenes"
	"github.com/zyebytevt/starburst-go/lib"
)

func setSceneCallback(button *lib.Button) error {
	button.SetHighlight(lib.HighlightInProgress)

	_, err := obsClient.Scenes.SetCurrentProgramScene(&scenes.SetCurrentProgramSceneParams{SceneName: button.Config.Parameters["scene_name"].(string)})
	return err
}

func toggleSourceVisibilityCallback(button *lib.Button) error {
	button.SetHighlight(lib.HighlightInProgress)

	sceneName := button.Config.Parameters["scene_name"].(string)

	sceneItemId, err := getSceneItemId(sceneName, button.Config.Parameters["source_name"].(string))

	if err != nil {
		return err
	}

	itemVisible, err := getSceneItemVisibility(sceneName, sceneItemId)

	if err != nil {
		return err
	}

	itemVisible = !itemVisible

	setSceneItemVisibility(sceneName, sceneItemId, itemVisible)

	if itemVisible {
		button.SetHighlight(lib.HighlightActive)
	} else {
		button.SetHighlight(lib.HighlightNone)
	}

	return nil
}
