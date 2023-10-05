.PHONY: build
build:
			go build -v ./cmd/app

.PHONY: test
test:
			go test -v -race ./...

.PHONY: compose-build
compose-build:
	docker compose build

.PHONY: compose-up
compose-up:
	docker compose up

.DEFAULT_GOAL := build