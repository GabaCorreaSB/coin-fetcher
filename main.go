//
// Copyright (c) 2023 Gabriel Correa <gabriel.correasb@protonmail.com>
//

package main

import (
	"flag"
)

func main() {

	listenAddr := flag.String("listenaddr", ":9899", "listen address the service is running")
	flag.Parse()

	svc := NewLogService(NewMetricService(&priceFetcher{}))

	server := NewJSONAPIServer(*listenAddr, svc)
	server.Run()
}
