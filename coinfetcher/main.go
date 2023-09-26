package main

import (
	"flag"
	"net/http"

	// Importing services created for our API
	coinApi "coinfetcher/api"
	healthService "coinfetcher/services/health"
	logUtils "coinfetcher/services/log"
	metricsUtils "coinfetcher/services/metrics"
	priceService "coinfetcher/services/price"
)

func main() {
	// Define a command-line flag to specify the listening address.
	listenAddr := flag.String("listenaddr", ":9899", "listen address for the service to run")
	flag.Parse()

	// Create instances of the price service and health checker.
	priceFetcher := priceService.NewPriceFetcher()
	healthChecker := healthService.NewHealthChecker()

	// Create instances of log and metrics services for price and health.
	coinService := logUtils.NewPriceLogService(metricsUtils.NewPriceMetricService(priceFetcher))
	healthService := logUtils.NewHealthLogService(metricsUtils.NewHealthMetricService(healthChecker))

	// Create a JSON API server instance with the specified services and listening address.
	server := coinApi.NewJSONAPIServer(*listenAddr, coinService, healthService)

	// Serve the REDOC Swagger UI HTML.
	http.Handle("/swagger/redoc.html", http.FileServer(http.Dir("./docs")))

	// Serve the Fast API-like Swagger UI.
	http.Handle("/swagger/swagger.html", http.FileServer(http.Dir("./docs")))

	// Serve the Swagger JSON file as swagger.json.
	http.Handle("/swagger/swagger.json", http.FileServer(http.Dir("./docs")))

	// Uncomment the line below if you have a Swagger YAML file as swagger.yaml.
	// http.Handle("/swagger/swagger.yaml", http.FileServer(http.Dir("./docs")))

	// Start the API server.
	server.Run()
}
