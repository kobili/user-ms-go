build:
	go build -o app

run:
	./app

run-and-build: build run
rb: run-and-build
