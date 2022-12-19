package twitch

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/magicmonkey/go-streamdeck"
	"github.com/nicklaw5/helix/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/zyebytevt/starburst-go/lib"
)

var twitchClient *helix.Client
var actionCallbacks map[string]lib.ActionCallback = map[string]lib.ActionCallback{
	"set_marker": setMarkerCallback,
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
	}

	configs, err := lib.GetConfigsForKey("twitch.buttons")
	if err != nil {
		return err
	}

	for _, config := range configs {
		lib.CreateButtonFromConfig(streamDeck, config, actionCallbacks)
	}

	return nil
}

func HandleAuthCallback(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "You are now authenticated. Feel free to close this tab.")

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
