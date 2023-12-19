.PHONY: build
build:
	go build -v ./cmd/apiWork
.DEFAULT_GOAL := build