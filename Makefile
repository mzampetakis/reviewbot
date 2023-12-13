# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

## build: build the cmd/rest application
.PHONY: build
build:
	go build -o=/tmp/bin/reviewbot ./cmd/reviewbot

## run: run the cmd/rest application
.PHONY: run
run: build
	/tmp/bin/reviewbot

# ==================================================================================== #
# RUNNING
# ==================================================================================== #
## start: start all necessary containers and run the app
.PHONY: start
start:
	touch .env
	docker-compose stop app
	docker-compose rm --force app
	docker-compose up -d --build app

## logs: watch the logs of all containers
.PHONY: logs
logs:
	docker-compose logs -f

## stop: stop all application containers
.PHONY: stop
stop:
	docker-compose stop

## clean: remove all containers and volumes of the app
.PHONY: clean
clean:
	docker-compose rm --stop --force
	docker image rm 'reviewbot-app'

