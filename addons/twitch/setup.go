package twitch

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/magicmonkey/go-streamdeck"
	"github.com/nicklaw5/helix/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/zyebytevt/streaming-backend/lib"
)

var twitchClient *helix.Client
var actionCallbacks map[string]func(userData any) = map[string]func(userData any){
	"set-marker": func(userData any) {
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
	},
}

type ActionConfig struct {
	ActionName  string `mapstructure:"action_name"`
	ButtonText  string `mapstructure:"button_text"`
	ButtonIndex int    `mapstructure:"button_index"`
}

func Setup(streamDeck *streamdeck.StreamDeck) error {
	logrus.Info("Initializing Twitch addon...")

	var err error
	twitchClient, err = helix.NewClient(&helix.Options{
		ClientID:     viper.GetString("twitch.client_id"),
		ClientSecret: viper.GetString("twitch.client_secret"),
		RedirectURI:  "http://localhost:7001/auth-callback",
	})

	if err != nil {
		return err
	}

	if err = updateTokens(); err != nil {
		// refresh token outdated or missing, re-auth
		// now set up the auth URL
		scopes := []string{"user:edit:broadcast"}
		url := twitchClient.GetAuthorizationURL(&helix.AuthorizationURLParams{
			ResponseType: "code", // or "token"
			Scopes:       scopes,
			ForceVerify:  false,
		})

		logrus.Infof("Auth to Twitch with URL in browser: %s", url)
	} else {
		logrus.WithError(err).Error("Failed to update token.")
	}

	configs := make([]*ActionConfig, 0)
	viper.UnmarshalKey("twitch.actions", &configs)

	for _, config := range configs {
		callback, exists := actionCallbacks[config.ActionName]
		if !exists {
			logrus.Warningf("Twitch action callback '%s' does not exist.", config.ActionName)
		}

		lib.NewButton(streamDeck, config.ButtonIndex, config.ButtonText, callback, config)
	}

	http.HandleFunc("/auth-callback", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "authed, like truthed -Some wise woman 2020")

		code := r.URL.Query().Get("code")

		resp, err := twitchClient.RequestUserAccessToken(code)
		if err != nil {
			panic(err)
		}

		access_token := resp.Data.AccessToken
		// Set the access token on the client
		twitchClient.SetUserAccessToken(access_token)

		refresh_token := resp.Data.RefreshToken
		// Put the refresh token in a file for later
		ioutil.WriteFile("twitch_refresh_token", []byte(refresh_token), 0644)
	})

	return nil
}

func updateTokens() error {
	data, err := ioutil.ReadFile("twitch_refresh_token")
	if err != nil {
		return err
	}
	refreshToken := string(data)

	resp, err := twitchClient.RefreshUserAccessToken(refreshToken)
	if err != nil {
		return err
	}

	access_token := resp.Data.AccessToken
	refresh_token := resp.Data.RefreshToken
	ioutil.WriteFile("twitch_refresh_token", []byte(refresh_token), 0644)
	// Set the access token on the client
	twitchClient.SetUserAccessToken(access_token)
	logrus.Info("Twitch tokens updated")

	return nil
}
