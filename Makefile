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
