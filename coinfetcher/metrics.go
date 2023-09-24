//
// Copyright (c) 2023 Gabriel Correa <gabriel.correasb@protonmail.com>
//

package main

import (
	"context"
	"fmt"
	"time"

	healthService "coinfetcher/services/health"
	priceService "coinfetcher/services/price"
)

type metricPriceService struct {
	next priceService.PriceFetcher
}

type metricHealthService struct {
	next healthService.HealthChecker
}

func NewPriceMetricService(next priceService.PriceFetcher) priceService.PriceFetcher {
	return &metricPriceService{
		next: next,
	}
}

func NewHealthMetricService(next healthService.HealthChecker) healthService.HealthChecker {
	return &metricHealthService{
		next: next,
	}
}

// FetchPrice fetches cryptocurrency price and logs metrics.
func (s *metricPriceService) FetchPrice(ctx context.Context, ticker string) (price float64, vol24Hr float64, timestamp time.Time, err error) {
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

func (s *metricHealthService) CheckHealth(ctx context.Context) (status string, geckoStatus string, timestamp time.Time, err error) {
	status, geckoStatus, timestamp, err = s.next.CheckHealth(ctx)
	if err != nil {
		fmt.Printf("Error getting status for Gecko API: %s\n", err)
	} else {
		fmt.Printf("Successfully fetched status for Gecko API")
		fmt.Printf("Status: %s\n", status)
		fmt.Printf("Gecko Status: %s\n", geckoStatus)
		fmt.Printf("Check Status Timestamp: %s\n", timestamp.String())
	}
	return status, geckoStatus, timestamp, err

}
