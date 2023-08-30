.PHONY: build
build:
	go build -v main.go

.DEFAULT_GOAL := build