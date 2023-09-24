//
// Copyright (c) 2023 Gabriel Correa <gabriel.correasb@protonmail.com>
//

package main

import (
	"flag"
	"net/http"

	// Importing services created to our api
	coinApi "coinfetcher/api"
	healthService "coinfetcher/services/health"
	logUtils "coinfetcher/services/log"
	metricsUtils "coinfetcher/services/metrics"
	priceService "coinfetcher/services/price"
)

func main() {
	listenAddr := flag.String("listenaddr", ":9899", "listen address the service is running")
	flag.Parse()

	priceFetcher := priceService.NewPriceFetcher()
	healthChecker := healthService.NewHealthChecker()

	coinService := logUtils.NewPriceLogService(metricsUtils.NewPriceMetricService(priceFetcher))
	healthService := logUtils.NewHealthLogService(metricsUtils.NewHealthMetricService(healthChecker))

	server := coinApi.NewJSONAPIServer(*listenAddr, coinService, healthService)

	// Serve the REDOC Swagger UI HTML
	http.Handle("/swagger/redoc.html", http.FileServer(http.Dir("./docs")))

	// Serve the Fast API like swagger UI
	http.Handle("/swagger/swagger.html", http.FileServer(http.Dir("./docs")))
	// Serve the Swagger JSON file as swagger.json
	http.Handle("/swagger/swagger.json", http.FileServer(http.Dir("./docs")))
	// Uncomment the line below if you have a Swagger YAML file as swagger.yaml
	http.Handle("/swagger/swagger.yaml", http.FileServer(http.Dir("./docs")))

	server.Run()
}
