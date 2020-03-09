SHELL := /bin/bash
APP_NAME = detection-api
VERSION ?= 0.1
SUSPICIOUS_SPEED ?= 100
NUM_OF_EVENTS ?= 3000

.PHONY: clean dependencies build test run run-image build-image clean run-generator

clean:
	rm $(PWD)/resources/event-db/event_db.db || true

dependencies:
	go mod vendor; go mod tidy

build: dependencies test
	go build -o $(PWD)/bin/$(APP_NAME); chmod +x $(PWD)/bin/$(APP_NAME)

test: dependencies
	go test -coverprofile=cover.out ./... -v

run: clean build
	SUSPICIOUS_SPEED=$(SUSPICIOUS_SPEED) $(PWD)/bin/$(APP_NAME)

build-image: clean build
	docker build --no-cache -t frankiennamdi/detection-api:$(VERSION) .

run-generator:
	go run generator/event_generator.go -num=$(NUM_OF_EVENTS)

docker-run:
	docker stop detection-api || true; docker rm detection-api || true;\
	docker run -e CONFIG_FILE=/app/resources/config.yml \
	-e DB_MIGRATION_LOC=/app/migrations \
	-e IP_GEO_DB_LOC=/app/resources/geo-database/GeoLite2-City.mmdb \
	-e SUSPICIOUS_SPEED=$(SUSPICIOUS_SPEED) \
	--name detection-api -p 3000:3000 frankiennamdi/detection-api:$(VERSION)

run-image: build-image docker-run

restart-image: docker-run
