package types

import "time"

type PriceResponse struct {
	Ticker    string    `json:"ticker"`
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
	Vol24Hr   float64   `json:"vol24Hr"`
}

type HealthResponse struct {
	Status         string    `json:"status"`
	GeckoApiStatus string    `json:"geckoapistatus"`
	Timestamp      time.Time `json:"timestamp"`
}
