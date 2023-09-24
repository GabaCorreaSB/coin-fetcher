//
// Copyright (c) 2023 Gabriel Correa <gabriel.correasb@protonmail.com>
//

package main

import (
	"context"
	"time"

	priceService "coinfetcher/services"

	log "github.com/sirupsen/logrus"
)

type logService struct {
	next priceService.PriceFetcher
}

func NewLogService(next priceService.PriceFetcher) priceService.PriceFetcher {
	return &logService{
		next: next,
	}
}

func (s *logService) FetchPrice(ctx context.Context, ticker string) (price float64, vol24Hr float64, timestamp time.Time, err error) {
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
