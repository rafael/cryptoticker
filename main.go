package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/rafael/cryptoticker/util"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/sheets/v4"
)

type CoinTicker struct {
	Id                             string      `json:"id"`
	Name                           string      `json:"name"`
	Symbol                         string      `json:"symbol"`
	Rank                           json.Number `json:"rank,Number"`
	PriceUsd                       json.Number `json:"price_usd"`
	TwentyFourHVolumeUsd           json.Number `json:"24h_volume_usd"`
	MarketCapUsd                   json.Number `json:"market_cap_usd"`
	AvailableSupply                json.Number `json:"available_supply"`
	TotalSupply                    json.Number `json:"total_supply"`
	MaxSupply                      json.Number `json:"max_supply"`
	PercentageChangeOneHour        json.Number `json:"percent_change_1h"`
	PercentageChangeTwentyFourHour json.Number `json:"percent_change_24h"`
	PercentageChangeSevenDays      json.Number `json:"percent_change_7d"`
}

func main() {
	config := util.GetConfig()
	tickers, err := getTicker()
	if err != nil {
		log.Fatal(err)
	}
	targetCoins := make(map[string]CoinTicker)
	for _, targetCoin := range config.TargetCoins {
		targetCoins[targetCoin] = CoinTicker{}
	}
	for _, ticker := range tickers {
		if _, ok := targetCoins[ticker.Id]; ok {
			targetCoins[ticker.Id] = ticker
		}
	}
	oauthToken, err := util.ParseOauthToken()
	if err != nil {
		log.Fatal(err)
	}
	updateGoogleSheet(oauthToken, targetCoins)
}

func getTicker() (tickers []CoinTicker, err error) {
	url := "https://api.coinmarketcap.com/v1/ticker"

	spaceClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		return nil, err
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, err
	}

	ticker := []CoinTicker{}
	jsonErr := json.Unmarshal(body, &ticker)
	if jsonErr != nil {
		return nil, err
	}
	return ticker, nil
}

func updateGoogleSheet(oauthToken *oauth2.Token, tickers map[string]CoinTicker) {
	ctx := context.Background()
	config, err := util.GetOauth2Config()
	if err != nil {
		log.Fatalf("Unable to get oauth config %v", err)
	}
	client := config.Client(ctx, oauthToken)
	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets Client %v", err)
	}

	spreadsheetId := util.GetConfig().SpreadSheetId
	writeRange := "coin_overview!A1"

	var vr sheets.ValueRange

	updatedAt := fmt.Sprintf("Last Updated: %s", time.Now().Format(time.RFC1123))
	myval := []interface{}{updatedAt}
	vr.Values = append(vr.Values, myval)

	_, err = srv.Spreadsheets.Values.Update(spreadsheetId, writeRange, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet. %v", err)
	}

	writeRange = "coin_overview!A4:F"
	vr = sheets.ValueRange{}
	for _, coin := range tickers {
		myval := []interface{}{coin.Name, coin.PriceUsd, coin.MarketCapUsd, coin.TwentyFourHVolumeUsd, coin.PercentageChangeSevenDays, coin.Rank}
		vr.Values = append(vr.Values, myval)
	}
	_, err = srv.Spreadsheets.Values.Update(spreadsheetId, writeRange, &vr).ValueInputOption("RAW").Do()
}
