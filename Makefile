build:
	go build -o bin/pricefetcher

run: build
	./bin/pricefetcher

swag:
	${HOME}/go/bin/swag init