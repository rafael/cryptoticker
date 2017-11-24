package cryptoticker

import (
	"encoding/json"
	"log"

	b64 "encoding/base64"

	"github.com/kelseyhightower/envconfig"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var config Config

type Config struct {
	TargetCoins        []string `envconfig:"target_coins"`
	OauthToken         string   `envconfig:"google_oauth_token"`
	GoogleClientSecret string   `envconfig:"google_client_secret"`
	SpreadSheetId      string   `envconfig:"google_spreadsheet_id"`
}

func GetConfig() Config {
	return config
}

func GetOauth2Config() (*oauth2.Config, error) {
	encodedSecret := config.GoogleClientSecret
	decodedSecret, err := b64.StdEncoding.DecodeString(encodedSecret)
	if err != nil {
		return nil, err
	}
	return google.ConfigFromJSON(decodedSecret, "https://www.googleapis.com/auth/spreadsheets")
}

func ParseOauthToken() (token *oauth2.Token, err error) {
	encodedToken := config.OauthToken
	decodedToken, err := b64.StdEncoding.DecodeString(encodedToken)
	if err != nil {
		return nil, err
	}
	var oauthToken oauth2.Token
	err = json.Unmarshal(decodedToken, &oauthToken)
	if err != nil {
		return nil, err
	}
	return &oauthToken, nil
}

func init() {
	err := envconfig.Process("crypto", &config)
	if err != nil {
		log.Fatal(err.Error())
	}
}
