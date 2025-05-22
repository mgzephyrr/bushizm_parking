GO := go

APP_NAME := app.exe
PATH_MAIN := $(CURDIR)/cmd
PATH_BIN := $(CURDIR)/bin

deps:
	$(GO) get -u ./...
	$(GO) mod tidy
	$(GO) mod vendor

build: deps
	$(GO) build -o $(PATH_BIN)/$(APP_NAME) $(PATH_MAIN)

run: build
	cd $(PATH_MAIN) && $(PATH_BIN)/$(APP_NAME)

docker-build: deps
	$(GO) build -o $(PATH_BIN)/$(APP_NAME) $(PATH_MAIN)

docker-run: docker-build
	cd $(PATH_MAIN) && $(PATH_BIN)/$(APP_NAME)
