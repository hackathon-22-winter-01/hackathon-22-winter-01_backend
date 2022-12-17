SHELL := /bin/bash
APP_NAME := app

# handle all targets as .PHONY
CMDS=$(shell grep -E -o "^[a-z_-]+" ./Makefile)
.PHONY: $(CMDS)

all: clean mod build

clean:
	rm -rf bin/$(APP_NAME)

mod:
	go mod download

build:
	go build -o bin/$(APP_NAME)

run:
	go run main.go

lint:
	golangci-lint run --fix ./...

go-gen:
	go generate ./...

TESTFLAGS ?= -cover -shuffle=on -race -timeout=30s
test:
	go test -v ./... $(TESTFLAGS)

FRONTEND_PATH := ./bin/frontend
update-frontend:
	[ -d $(FRONTEND_PATH) ] || git clone git@github.com:hackathon-22-winter-01/front-end.git $(FRONTEND_PATH)
	git pull origin main
	yarn --cwd $(FRONTEND_PATH) install
	yarn --cwd $(FRONTEND_PATH) build
	go run main.go
