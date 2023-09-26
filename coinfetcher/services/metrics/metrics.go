package metrics_utils

import (
	"context"
	"fmt"
	"time"

	// Importing service packages for health and price data.
	healthService "coinfetcher/services/health"
	priceService "coinfetcher/services/price"
)

// Definition of the metricPriceService struct, which extends priceService.PriceFetcher.
type metricPriceService struct {
	next priceService.PriceFetcher // The 'next' field holds an instance of the underlying price service.
}

// Definition of the metricHealthService struct, which extends healthService.HealthChecker.
type metricHealthService struct {
	next healthService.HealthChecker // The 'next' field holds an instance of the underlying health service.
}

// Factory function to create a new metricPriceService instance.
// It accepts the underlying price service as a parameter and returns a priceService.PriceFetcher.
func NewPriceMetricService(next priceService.PriceFetcher) priceService.PriceFetcher {
	return &metricPriceService{
		next: next,
	}
}

// Factory function to create a new metricHealthService instance.
// It accepts the underlying health service as a parameter and returns a healthService.HealthChecker.
func NewHealthMetricService(next healthService.HealthChecker) healthService.HealthChecker {
	return &metricHealthService{
		next: next,
	}
}

// FetchPrice method of metricPriceService.
// It fetches cryptocurrency price and logs metrics, delegating the actual fetching to the underlying service.
func (s *metricPriceService) FetchPrice(ctx context.Context, ticker string) (price float64, vol24Hr float64, timestamp time.Time, err error) {
	price, vol24Hr, timestamp, err = s.next.FetchPrice(ctx, ticker) // Delegates the fetching to the underlying service.
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

// CheckHealth method of metricHealthService.
// It checks the health of a service and logs metrics, delegating the actual check to the underlying service.
func (s *metricHealthService) CheckHealth(ctx context.Context) (status string, geckoStatus string, timestamp time.Time, err error) {
	status, geckoStatus, timestamp, err = s.next.CheckHealth(ctx) // Delegates the health check to the underlying service.
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
