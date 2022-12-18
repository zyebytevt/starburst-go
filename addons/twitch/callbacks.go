package twitch

import (
	"github.com/nicklaw5/helix/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/zyebytevt/streaming-backend/lib"
)

func setMarkerCallback(config *lib.ButtonConfig) {
	isValid, _, _ := twitchClient.ValidateToken(twitchClient.GetUserAccessToken())
	if !isValid {
		updateTokens()
	}

	resp_mark, err := twitchClient.CreateStreamMarker(&helix.CreateStreamMarkerParams{
		UserID:      viper.GetString("twitch.user_id"),
		Description: "Streamdeck marks the spot",
	})

	if err != nil {
		logrus.WithError(err).Error("Twitch Helix request failed.")
		return
	}

	if resp_mark.Error != "" {
		logrus.WithField("TwitchError", resp_mark.ErrorMessage).Error("Failed to set stream marker.")
		return
	}

	logrus.Infof("Created stream marker at %v.", resp_mark.Data.CreateStreamMarkers[0].CreatedAt)
}
