package obs

import "github.com/andreykaipov/goobs/api/requests/sceneitems"

func getSceneItemId(sceneName string, sourceName string) (float64, error) {
	sceneItemId, err := obsClient.SceneItems.GetSceneItemId(&sceneitems.GetSceneItemIdParams{
		SceneName:  sceneName,
		SourceName: sourceName,
	})

	if err != nil {
		return 0, err
	}

	return sceneItemId.SceneItemId, nil
}

func getSceneItemVisibility(sceneName string, sceneItemId float64) (bool, error) {
	itemEnabled, err := obsClient.SceneItems.GetSceneItemEnabled(&sceneitems.GetSceneItemEnabledParams{
		SceneName:   sceneName,
		SceneItemId: sceneItemId,
	})

	if err != nil {
		return false, err
	}

	return itemEnabled.SceneItemEnabled, nil
}

func setSceneItemVisibility(sceneName string, sceneItemId float64, visible bool) error {
	_, err := obsClient.SceneItems.SetSceneItemEnabled(&sceneitems.SetSceneItemEnabledParams{
		SceneName:        sceneName,
		SceneItemId:      sceneItemId,
		SceneItemEnabled: &visible,
	})

	return err
}
