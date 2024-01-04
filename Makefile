.PHONY: build
build:
	go build -v ./cmd/server
	go build -v ./cmd/publisher
	cd ./docker/ && docker-compose up
	
.PHONY: test
test:
	cd ./docker/ && docker-compose up
	go test -v -race -timeout 30s ./...
	cd ./docker/ && docker-compose down

.PHONY: clean
clean:
	rm -f ./publisher
	rm -f ./server
	cd ./docker/ && docker-compose down

.DEFAULT_GOAL := build