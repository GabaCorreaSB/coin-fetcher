//
// Copyright (c) 2023 Gabriel Correa <gabriel.correasb@protonmail.com>
//

package health_service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Creating structure for status data.
var data struct {
	GeckoStatus string `json:"gecko_says"`
}

// healthChecker is an interface that can fetch a price.
type HealthChecker interface {
	CheckHealth(context.Context) (string, string, time.Time, error)
}

// healtChecker implements the HealthChecker interface.
type healthChecker struct{}

func NewHealthChecker() HealthChecker {
	return &healthChecker{}
}

// Fetching the data
func (s *healthChecker) CheckHealth(ctx context.Context) (string, string, time.Time, error) {
	price, vol24Hr, timestamp, err := CheckGeckoHealth()
	if err != nil {
		return "", "", time.Time{}, fmt.Errorf("failed to fetch crypto price: %v", err)
	}
	return price, vol24Hr, timestamp, nil
}

// Function that takes care of checking gecko api status
func CheckGeckoHealth() (string, string, time.Time, error) {
	const coingeckoAPI = "https://api.coingecko.com/api/v3/ping"

	client := &http.Client{Timeout: 10 * time.Second} // Add a timeout for the HTTP client.

	// Make a GET request to the CoinGecko API.
	resp, err := client.Get(coingeckoAPI)
	if err != nil {
		return "", "", time.Time{}, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", time.Time{}, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&data); err != nil {
		return "", "", time.Time{}, fmt.Errorf("failed to decode JSON response: %v", err)
	}

	healthData := data

	unixTimeNow := time.Now().UTC().Unix()
	geckoStatus := healthData.GeckoStatus

	status := "Not Running"
	if geckoStatus != "" {
		status = "Running"
	}

	timestamp := time.Unix(unixTimeNow, 0)

	return status, geckoStatus, timestamp, nil
}
