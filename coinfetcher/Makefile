build:
	go build -o bin/pricefetcher

run: build
	mkdir logs
	./bin/pricefetcher >> logs/coin_api.log

swag:
	${HOME}/go/bin/swag init