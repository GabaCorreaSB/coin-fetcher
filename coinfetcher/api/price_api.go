// This section specifies the package name and imports required for the code.
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

// APIFunc is a type representing a function that handles API requests.
type APIFunc func(context.Context, http.ResponseWriter, *http.Request) error

// PriceResponse represents the response format for price-related endpoints.
type PriceResponse struct {
	Ticker    string    `json:"ticker"`
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
	Vol24Hr   float64   `json:"vol24Hr"`
}

// HealthResponse represents the response format for health-related endpoints.
type HealthResponse struct {
	Status         string    `json:"status"`
	GeckoApiStatus string    `json:"geckoapistatus"`
	Timestamp      time.Time `json:"timestamp"`
}

// JSONAPIServer represents a JSON API server.
type JSONAPIServer struct {
	listenAddr     string
	pricingService priceService.PriceFetcher
	statusService  healthService.HealthChecker
}

// NewJSONAPIServer creates a new instance of the JSONAPIServer.
func NewJSONAPIServer(listenAddr string, pricingService priceService.PriceFetcher, statusService healthService.HealthChecker) *JSONAPIServer {
	return &JSONAPIServer{
		listenAddr:     listenAddr,
		pricingService: pricingService,
		statusService:  statusService,
	}
}

// Run starts the JSON API server.
func (s *JSONAPIServer) Run() {
	http.HandleFunc("/v1/price", s.makeHTTPHandlerFunc(s.handleFetchPrice))
	http.HandleFunc("/v1/health", s.makeHTTPHandlerFunc(s.handleApiHealth))

	http.ListenAndServe(s.listenAddr, nil)
}

// makeHTTPHandlerFunc is a helper function to create an HTTP handler function.
func (s *JSONAPIServer) makeHTTPHandlerFunc(apiFn APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "requestID", rand.Intn(10000000))

		if err := apiFn(ctx, w, r); err != nil {
			s.writeJSON(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		}
	}
}

// handleFetchPrice handles the "Fetch coin price" endpoint.
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

// handleApiHealth handles the "Get Gecko API health status" endpoint.
func (s *JSONAPIServer) handleApiHealth(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	health, geckoStatus, timestamp, err := s.statusService.CheckHealth(ctx)
	if err != nil {
		return err
	}

	healthResponse := types.HealthResponse{
		Status:         health,
		GeckoApiStatus: geckoStatus,
		Timestamp:      timestamp,
	}

	return s.writeJSON(w, http.StatusOK, &healthResponse)
}

// writeJSON writes JSON responses with the specified status code.
func (s *JSONAPIServer) writeJSON(w http.ResponseWriter, statusCode int, v interface{}) error {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
