package log_utils

import (
	"context"
	"time"

	// Importing service packages for health and price data.
	healthService "coinfetcher/services/health"
	priceService "coinfetcher/services/price"

	log "github.com/sirupsen/logrus" // Importing the logrus package for logging.
)

// Definition of the logPriceService struct, which extends priceService.PriceFetcher.
type logPriceService struct {
	next priceService.PriceFetcher // The 'next' field holds an instance of the underlying price service.
}

// Definition of the logHealthService struct, which extends healthService.HealthChecker.
type logHealthService struct {
	next healthService.HealthChecker // The 'next' field holds an instance of the underlying health service.
}

// Factory function to create a new logPriceService instance.
// It accepts the underlying price service as a parameter and returns a priceService.PriceFetcher.
func NewPriceLogService(next priceService.PriceFetcher) priceService.PriceFetcher {
	return &logPriceService{
		next: next,
	}
}

// Factory function to create a new logHealthService instance.
// It accepts the underlying health service as a parameter and returns a healthService.HealthChecker.
func NewHealthLogService(next healthService.HealthChecker) healthService.HealthChecker {
	return &logHealthService{
		next: next,
	}
}

// FetchPrice method of logPriceService.
// It fetches cryptocurrency price, logs metrics, and adds log entries with relevant information.
func (s *logPriceService) FetchPrice(ctx context.Context, ticker string) (price float64, vol24Hr float64, timestamp time.Time, err error) {
	begin := time.Now() // Record the start time.

	// Delegate the price fetching to the underlying service.
	price, vol24Hr, timestamp, err = s.next.FetchPrice(ctx, ticker)

	// Create log fields to store relevant information.
	fields := log.Fields{
		"requestID": ctx.Value("requestID"), // Context value, if available.
		"took":      time.Since(begin),      // Time taken for the operation.
		"err":       err,                    // Error, if any.
		"price":     price,                  // Fetched price.
		"vol24Hr":   vol24Hr,                // 24-hour volume.
		"timestamp": timestamp,              // Price timestamp.
	}

	// Log the information using logrus with the "fetchPrice" log message.
	log.WithFields(fields).Info("fetchPrice")

	return price, vol24Hr, timestamp, err
}

// CheckHealth method of logHealthService.
// It checks the health of a service, logs metrics, and adds log entries with relevant information.
func (s *logHealthService) CheckHealth(ctx context.Context) (status string, geckoStatus string, timestamp time.Time, err error) {
	begin := time.Now() // Record the start time.

	// Delegate the health check to the underlying service.
	status, geckoStatus, timestamp, err = s.next.CheckHealth(ctx)

	// Create log fields to store relevant information.
	fields := log.Fields{
		"requestID":   ctx.Value("requestID"), // Context value, if available.
		"took":        time.Since(begin),      // Time taken for the operation.
		"err":         err,                    // Error, if any.
		"status":      status,                 // Service status.
		"geckoStatus": geckoStatus,            // Gecko API status.
		"timestamp":   timestamp,              // Check timestamp.
	}

	// Log the information using logrus with the "checkHealth" log message.
	log.WithFields(fields).Info("checkHealth")

	return status, geckoStatus, timestamp, err
}
