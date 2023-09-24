//
// Copyright (c) 2023 Gabriel Correa <gabriel.correasb@protonmail.com>
//

package main

import (
	"context"
	"fmt"
	"time"
)

type metricService struct {
	next PriceFetcher
}

func NewMetricService(next PriceFetcher) PriceFetcher {
	return &metricService{
		next: next,
	}
}

// FetchPrice fetches cryptocurrency price and logs metrics.
func (s *metricService) FetchPrice(ctx context.Context, ticker string) (price float64, vol24Hr float64, timestamp time.Time, err error) {
	price, vol24Hr, timestamp, err = s.next.FetchPrice(ctx, ticker)
	if err != nil {
		fmt.Printf("Error fetching price for ticker %s: %v\n", ticker, err)
	} else {
		fmt.Printf("Successfully fetched price for ticker %s:\n", ticker)
		fmt.Printf("Price: %f\n", price)
		fmt.Printf("24-Hour Volume: %f\n", vol24Hr)
		fmt.Printf("Price Timestamp: %s\n", timestamp.String())
	}
	return price, vol24Hr, timestamp, err
}
