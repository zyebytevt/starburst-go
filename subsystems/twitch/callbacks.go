package twitch

import (
	"errors"

	"github.com/nicklaw5/helix/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/zyebytevt/starburst-go/lib"
)

func setMarkerCallback(button *lib.Button) error {
	isValid, _, _ := twitchClient.ValidateToken(twitchClient.GetUserAccessToken())
	if !isValid {
		updateTokens()
	}

	resp_mark, err := twitchClient.CreateStreamMarker(&helix.CreateStreamMarkerParams{
		UserID:      viper.GetString("twitch.user_id"),
		Description: "Streamdeck marks the spot",
	})

	if err != nil {
		return err
	}

	if resp_mark.Error != "" {
		return errors.New(resp_mark.ErrorMessage)
	}

	button.FlashNotify(lib.NotifySuccess)
	logrus.Infof("Created stream marker at %v.", resp_mark.Data.CreateStreamMarkers[0].CreatedAt)
	return nil
}
