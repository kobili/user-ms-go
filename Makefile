build:
	cd go && go build -o ../app;

run:
	./app

run-and-build: build run
rb: run-and-build

dependencies:
	docker compose -f docker-compose.yml up
deps: dependencies

make-migrations:
	docker compose -f docker-compose.yml -f compose.migrate.yml run --rm migrate create -ext .sql -dir ./migrations -seq ${migration_name}
mm: make-migrations

migrate:
	docker compose -f docker-compose.yml -f compose.migrate.yml run --rm migrate -path /migrations -database postgres://postgres:password@postgres:5432/userms?sslmode=disable up
m: migrate

count ?=
migrate-down:
	docker compose -f docker-compose.yml -f compose.migrate.yml run --rm migrate -path /migrations -database postgres://postgres:password@postgres:5432/userms?sslmode=disable down ${count}
md: migrate-down

go-tidy:
	cd go && go mod tidy
