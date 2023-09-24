build:
	go build -o bin/coinfetcher

run: build
	mkdir logs
	./bin/coinfetcher >> logs/coin_api.log

swag:
	${HOME}/go/bin/swag init