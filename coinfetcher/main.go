//
// Copyright (c) 2023 Gabriel Correa <gabriel.correasb@protonmail.com>
//

package main

import (
	"flag"
	"net/http"

	coinApi "coinfetcher/api"
	healthService "coinfetcher/services/health"
	priceService "coinfetcher/services/price"
)

func main() {
	listenAddr := flag.String("listenaddr", ":9899", "listen address the service is running")
	flag.Parse()

	priceFetcher := priceService.NewPriceFetcher()
	healthChecker := healthService.NewHealthChecker()

	coinService := NewPriceLogService(NewPriceMetricService(priceFetcher))
	healthService := NewHealthLogService(NewHealthMetricService(healthChecker))

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
