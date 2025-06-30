build:
	go build -o app

run:
	./app

run-and-build: build run
rb: run-and-build

dependencies:
	docker compose -f docker-compose.yml up
deps: dependencies
