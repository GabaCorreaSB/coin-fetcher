package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Creating structure for data to be fetched.
var data map[string]struct {
	USD           float64 `json:"usd"`
	LastUpdatedAt int64   `json:"last_updated_at"`
	Vol24Hr       float64 `json:"usd_24h_vol"`
}

// PriceFetcher is an interface that can fetch a price.
type PriceFetcher interface {
	FetchPrice(context.Context, string) (float64, float64, time.Time, error)
}

// priceFetcher implements the PriceFetcher interface.
type priceFetcher struct{}

func NewPriceFetcher() PriceFetcher {
	return &priceFetcher{}
}

// Fetching the data
func (s *priceFetcher) FetchPrice(ctx context.Context, ticker string) (float64, float64, time.Time, error) {
	price, vol24Hr, timestamp, err := FetchCryptoPrice(ticker)
	if err != nil {
		return 0, 0, time.Time{}, fmt.Errorf("failed to fetch crypto price: %v", err)
	}
	return price, vol24Hr, timestamp, nil
}

// Function that takes care of creating our platform info from CoinGecko
func FetchCryptoPrice(ticker string) (float64, float64, time.Time, error) {
	const coingeckoAPI = "https://api.coingecko.com/api/v3/simple/price"

	client := &http.Client{Timeout: 10 * time.Second} // Add a timeout for the HTTP client.

	// Make a GET request to the CoinGecko API.
	resp, err := client.Get(fmt.Sprintf("%s?ids=%s&vs_currencies=usd&include_24hr_vol=true&include_last_updated_at=true", coingeckoAPI, ticker))
	if err != nil {
		return 0, 0, time.Time{}, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, time.Time{}, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&data); err != nil {
		return 0, 0, time.Time{}, fmt.Errorf("failed to decode JSON response: %v", err)
	}

	priceData, ok := data[ticker]
	if !ok {
		return 0, 0, time.Time{}, errors.New("could not find data for ticker")
	}

	price := priceData.USD
	vol24Hr := priceData.Vol24Hr
	timestamp := time.Unix(priceData.LastUpdatedAt, 0)

	return price, vol24Hr, timestamp, nil
}
