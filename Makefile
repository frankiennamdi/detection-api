SHELL := /bin/bash
APP_NAME = detection-api
VERSION ?= 0.1

dependencies:
	go mod vendor; go mod tidy

build: dependencies test
	go build -o $(PWD)/bin/$(APP_NAME); chmod +x $(PWD)/bin/$(APP_NAME)

test: dependencies
	go test -coverprofile=cover.out ./... -v

run: build
	$(PWD)/bin/$(APP_NAME)

build-image: build
	docker build --no-cache -t frankiennamdi/detection-api:$(VERSION) .

run-image: build-image
	docker stop detection-api || true; docker rm detection-api || true;\
	docker run -e "CONFIG_FILE=/app/resources/config.yml" \
	-e "DB_MIGRATION_LOC=/app/migrations" \
	-e "IP_GEO_DB_LOC=/app/resources/geo-database/GeoLite2-City.mmdb" \
	--name detection-api -p 3000:3000 frankiennamdi/detection-api:$(VERSION)