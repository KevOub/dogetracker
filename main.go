package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Reads config.json
type config struct {
	Username       string
	Coins          string
	DogeToken      string
	NomicsAPI      string
	DiscordWebhook string
	Intervals      int
	Thresholds     float64
}

// Thank goodness for https://mholt.github.io/json-to-go/
type NomicsAPI []struct {
	ID                string    `json:"id"`
	Currency          string    `json:"currency"`
	Symbol            string    `json:"symbol"`
	Name              string    `json:"name"`
	LogoURL           string    `json:"logo_url"`
	Status            string    `json:"status"`
	Price             string    `json:"price"`
	PriceDate         time.Time `json:"price_date"`
	PriceTimestamp    time.Time `json:"price_timestamp"`
	CirculatingSupply string    `json:"circulating_supply"`
	MarketCap         string    `json:"market_cap"`
	NumExchanges      string    `json:"num_exchanges"`
	NumPairs          string    `json:"num_pairs"`
	NumPairsUnmapped  string    `json:"num_pairs_unmapped"`
	FirstCandle       time.Time `json:"first_candle"`
	FirstTrade        time.Time `json:"first_trade"`
	FirstOrderBook    time.Time `json:"first_order_book"`
	Rank              string    `json:"rank"`
	RankDelta         string    `json:"rank_delta"`
	High              string    `json:"high"`
	HighTimestamp     time.Time `json:"high_timestamp"`
	OneD              struct {
		Volume             string `json:"volume"`
		PriceChange        string `json:"price_change"`
		PriceChangePct     string `json:"price_change_pct"`
		VolumeChange       string `json:"volume_change"`
		VolumeChangePct    string `json:"volume_change_pct"`
		MarketCapChange    string `json:"market_cap_change"`
		MarketCapChangePct string `json:"market_cap_change_pct"`
	} `json:"1d"`
	Three0D struct {
		Volume             string `json:"volume"`
		PriceChange        string `json:"price_change"`
		PriceChangePct     string `json:"price_change_pct"`
		VolumeChange       string `json:"volume_change"`
		VolumeChangePct    string `json:"volume_change_pct"`
		MarketCapChange    string `json:"market_cap_change"`
		MarketCapChangePct string `json:"market_cap_change_pct"`
	} `json:"30d"`
}

func main() {
	// Read json from disk
	file, _ := ioutil.ReadFile("config.json")

	settings := config{}

	err := json.Unmarshal([]byte(file), &settings)
	if err != nil {
		log.Fatal(err)
	}

	localMaximum := float64(1 / 10000000000000)
	localMinimum := float64(9999999999999999)
	previousPrice := 0.0
	for range time.Tick(time.Second * time.Duration(settings.Intervals)) {

		currentPrice := GetAPIBody(settings)[0].Price
		currentPriceInt, _ := strconv.ParseFloat(currentPrice, 64)
		if currentPriceInt > localMaximum {
			/* 			if localMaximum/currentPriceInt >= settings.Thresholds {
			   				DiscordPing(settings, "🟩 Currency has increased past the threshold ")
			   			}
			*/localMaximum = currentPriceInt
		}

		if currentPriceInt < localMinimum {

			/* 			if currentPriceInt/localMinimum >= settings.Thresholds {
			   				DiscordPing(settings, "🟥 Currency has decreased past the threshold ")
			   			}
			*/
			localMinimum = currentPriceInt
		}
		if currentPriceInt-previousPrice > 0 {
			message := fmt.Sprintf("🟩 PRICE: %s", currentPrice)
			DiscordPing(settings, message)

		} else {
			message := fmt.Sprintf("🟥 PRICE: %s", currentPrice)
			DiscordPing(settings, message)

		}
		previousPrice = currentPriceInt

	}
}

// DiscordPing sends Discord ping
func DiscordPing(settings config, message string) {
	// Send username and "Content"/ message
	params := fmt.Sprintf(`{
		"username":"%s",
		"content": "%s"
	}`, settings.Username, message)
	requestBody := strings.NewReader(params)

	// Post with the webhook generated being used
	resp, err := http.Post(
		settings.DiscordWebhook, "application/json; charset=UTF-8", requestBody,
	)

	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()

}

// GetCoinAmount returns current price of users dogecoin [dumb test suite]
func GetCoinAmount(settings config) float64 {

	url := fmt.Sprintf("https://dogechain.info/chain/Dogecoin/q/addressbalance/%s", settings.DogeToken)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	output, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	strVersion := fmt.Sprintf("%s", output)
	floatVersion, _ := strconv.ParseFloat(strVersion, 10)
	return floatVersion
}

// GetAPIBody returns current price of dogecoin
func GetAPIBody(settings config) NomicsAPI {
	// adds API key to url
	url := fmt.Sprintf("https://api.nomics.com/v1/currencies/ticker?key=%s&ids=DOGE&interval=1h,30d&per-page=100&page=1", settings.NomicsAPI)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	// Reads output body
	output, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	// converts to JSON
	apiBody := NomicsAPI{}

	err = json.Unmarshal(output, &apiBody)
	if err != nil {
		log.Fatal(err)
	}

	return apiBody
}
