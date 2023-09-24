//
// Copyright (c) 2023 Gabriel Correa <gabriel.correasb@protonmail.com>
//

// @host localhost:9899
// @BasePath /v1
// @schemes http
package price_api

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	healthService "coinfetcher/services/health"
	priceService "coinfetcher/services/price"
	"coinfetcher/types"
)

type APIFunc func(context.Context, http.ResponseWriter, *http.Request) error

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

type JSONAPIServer struct {
	listenAddr     string
	pricingService priceService.PriceFetcher
	statusService  healthService.HealthChecker
}

func NewJSONAPIServer(listenAddr string, pricingService priceService.PriceFetcher, statusService healthService.HealthChecker) *JSONAPIServer {
	return &JSONAPIServer{
		listenAddr:     listenAddr,
		pricingService: pricingService,
		statusService:  statusService,
	}
}

func (s *JSONAPIServer) Run() {
	http.HandleFunc("/v1/price", s.makeHTTPHandlerFunc(s.handleFetchPrice))
	http.HandleFunc("/v1/health", s.makeHTTPHandlerFunc(s.handleApiHealth))

	http.ListenAndServe(s.listenAddr, nil)
}

func (s *JSONAPIServer) makeHTTPHandlerFunc(apiFn APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "requestID", rand.Intn(10000000))

		if err := apiFn(ctx, w, r); err != nil {
			s.writeJSON(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		}
	}
}

// @Summary Fetch coin price Endpoint
// @Description Fetches the price of a given coin ticker.
// @ID fetchPrice
// @Produce json
// @Param ticker query string true "Coin ticker symbol"
// @Success 200 {object} PriceResponse
// @Router /v1/price [get]
func (s *JSONAPIServer) handleFetchPrice(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ticker := r.URL.Query().Get("ticker")

	price, vol24Hr, timestamp, err := s.pricingService.FetchPrice(ctx, ticker)
	if err != nil {
		return err
	}

	priceResp := types.PriceResponse{
		Price:     price,
		Ticker:    ticker,
		Timestamp: timestamp,
		Vol24Hr:   vol24Hr,
	}

	return s.writeJSON(w, http.StatusOK, &priceResp)
}

// @Summary Get Gecko API health status Endpoint
// @Description This gets the Gecko API health status
// @ID checkHealth
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /v1/health [get]
func (s *JSONAPIServer) handleApiHealth(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	health, geckoStatus, timestamp, err := s.statusService.CheckHealth(ctx)
	if err != nil {
		return err
	}

	healtResponse := types.HealthResponse{
		Status:         health,
		GeckoApiStatus: geckoStatus,
		Timestamp:      timestamp,
	}

	return s.writeJSON(w, http.StatusOK, &healtResponse)
}

func (s *JSONAPIServer) writeJSON(w http.ResponseWriter, statusCode int, v interface{}) error {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
