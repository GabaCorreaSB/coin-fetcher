//
// Copyright (c) 2023 Gabriel Correa <gabriel.correasb@protonmail.com>
//

package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/GabaCorreaSB/coin-checker/types"
)

type APIFunc func(context.Context, http.ResponseWriter, *http.Request) error

type JSONAPIServer struct {
	listenAddr string
	svc        PriceFetcher
}

func NewJSONAPIServer(listenAddr string, svc PriceFetcher) *JSONAPIServer {
	return &JSONAPIServer{
		listenAddr: listenAddr,
		svc:        svc,
	}
}

func (s *JSONAPIServer) Run() {
	http.HandleFunc("/price", s.makeHTTPHandlerFunc(s.handleFetchPrice))
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

func (s *JSONAPIServer) handleFetchPrice(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ticker := r.URL.Query().Get("ticker")

	price, vol24Hr, timestamp, err := s.svc.FetchPrice(ctx, ticker)
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

func (s *JSONAPIServer) writeJSON(w http.ResponseWriter, statusCode int, v interface{}) error {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
