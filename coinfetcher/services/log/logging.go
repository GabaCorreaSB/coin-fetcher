//
// Copyright (c) 2023 Gabriel Correa <gabriel.correasb@protonmail.com>
//

package log_utils

import (
	"context"
	"time"

	healthService "coinfetcher/services/health"
	priceService "coinfetcher/services/price"

	log "github.com/sirupsen/logrus"
)

type logPriceService struct {
	next priceService.PriceFetcher
}

type logHealthService struct {
	next healthService.HealthChecker
}

func NewPriceLogService(next priceService.PriceFetcher) priceService.PriceFetcher {
	return &logPriceService{
		next: next,
	}
}

func NewHealthLogService(next healthService.HealthChecker) healthService.HealthChecker {
	return &logHealthService{
		next: next,
	}
}

func (s *logPriceService) FetchPrice(ctx context.Context, ticker string) (price float64, vol24Hr float64, timestamp time.Time, err error) {
	begin := time.Now()
	price, vol24Hr, timestamp, err = s.next.FetchPrice(ctx, ticker)

	fields := log.Fields{
		"requestID": ctx.Value("requestID"),
		"took":      time.Since(begin),
		"err":       err,
		"price":     price,
		"vol24Hr":   vol24Hr,
		"timestamp": timestamp,
	}

	log.WithFields(fields).Info("fetchPrice")

	return price, vol24Hr, timestamp, err
}

func (s *logHealthService) CheckHealth(ctx context.Context) (status string, geckoStatus string, timestamp time.Time, err error) {
	begin := time.Now()
	status, geckoStatus, timestamp, err = s.next.CheckHealth(ctx)

	fields := log.Fields{
		"requestID":   ctx.Value("requestID"),
		"took":        time.Since(begin),
		"err":         err,
		"status":      status,
		"geckoStatus": geckoStatus,
		"timestamp":   timestamp,
	}

	log.WithFields(fields).Info("checkHealth")

	return status, geckoStatus, timestamp, err
}
