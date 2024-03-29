build-redis:
	cd ./redis && go build -o bin/$(basename $(pwd))

build-elasticsearch:
	cd ./elasticsearch && go build -o bin/$(basename $(pwd))

build-jupyter:
	cd ./jupyter && go build -o bin/$(basename $(pwd))

build-graphql:
	cd ./graphql && go build -o bin/$(basename $(pwd))

test-all:
	go test ./redis
	go test ./elasticsearch
	go test ./jupyter
	go test ./graphql

build-all: build-redis build-elasticsearch build-jupyter build-graphql
