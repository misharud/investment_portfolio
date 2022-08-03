GO = $(shell which go)
PACKAGENAME := $(shell go list -m -f '{{.Path}}')
GITVERSION := $(shell git describe --dirty --always --tags --long)
DATE := $(shell date -u '+%Y-%m-%d-%H:%M UTC')
LDFLAGS = -ldflags " \
    -X '${PACKAGENAME}/internal/pkg/config.Version=${GITVERSION}' \
    -X '${PACKAGENAME}/internal/pkg/config.BuildDate=${DATE}' \
  " \

services-up:
	docker-compose -f ./deployment/local/docker-compose.services.yml up --build -d

services-down:
	docker-compose -f ./deployment/local/docker-compose.services.yml down

lint-ci:
	./tools.sh lint
