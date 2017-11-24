package main

import (
	"encoding/json"
	"fmt"
	"log"

	b64 "encoding/base64"
	"github.com/rafael/cryptoticker"
	"golang.org/x/oauth2"
)

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}
func printToken(token *oauth2.Token) {
	jsonToken, err := json.Marshal(token)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	sEnc := b64.StdEncoding.EncodeToString(jsonToken)
	fmt.Printf("Save this env: CRYPTO_GOOGLE_OAUTH_TOKEN=%s", sEnc)
}

func main() {
	config, err := cryptoticker.GetOauth2Config()
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	tok := getTokenFromWeb(config)
	printToken(tok)
}
