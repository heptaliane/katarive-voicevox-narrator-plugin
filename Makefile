TARGET := cpu

.PHONY: build test format run-voicevox

build:
	go build .

test:
	go test ./...

format:
	go fmt ./...

run-voicevox:
	docker compose -f docker-compose.$(TARGET).yaml up -d

stop-voicevox:
	docker compose -f docker-compose.$(TARGET).yaml down
