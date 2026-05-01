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

generate: run-voicevox
	mkdir -p gen/voicevox
	curl http://localhost:50021/openapi.json -o openapi.json
	sed -E 's/\{"anyOf":\[\{"type":"([^"]+)"\},\{"type":"null"\}\](\})?/{"type":"\1","nullable":true\2/g' openapi.json > openapi_v3_0.json
	oapi-codegen -package voicevox -generate types,client openapi_v3_0.json > gen/voicevox/voicevox.go
	rm openapi.json openapi_v3_0.json
