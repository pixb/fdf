BINARY_NAME := fdf
VERSION := $(shell cat VERSION 2>/dev/null | tr -d '[:space:]')
.PHONY: build install clean test

build:
	@echo "Building $(BINARY_NAME) version $(VERSION)"
	CGO_ENABLED=0 go build -o $(BINARY_NAME) -ldflags=-X=main.Version=$(VERSION) main.go

install:
	@echo "Installing $(BINARY_NAME) version $(VERSION)"
	CGO_ENABLED=0 go install -ldflags=-X=main.Version=$(VERSION)

clean:
	rm -f $(BINARY_NAME)

test:
	go test -v ./...
