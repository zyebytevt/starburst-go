package obs

import (
	"github.com/andreykaipov/goobs/api/requests/scenes"
	"github.com/zyebytevt/starburst-go/lib"
)

func setSceneCallback(button *lib.Button) error {
	button.SetDecorator(lib.InProgressStateDecorator)

	_, err := obsClient.Scenes.SetCurrentProgramScene(&scenes.SetCurrentProgramSceneParams{SceneName: button.Config.Parameters["scene_name"].(string)})
	return err
}

func toggleSourceVisibilityCallback(button *lib.Button) error {
	button.SetDecorator(lib.InProgressStateDecorator)

	sceneName := button.Config.Parameters["scene_name"].(string)
	sceneItemId := button.UserData["scene_item_id"].(float64)

	itemVisible, err := getSceneItemVisibility(sceneName, sceneItemId)

	if err != nil {
		return err
	}

	itemVisible = !itemVisible

	return setSceneItemVisibility(sceneName, sceneItemId, itemVisible)
}
