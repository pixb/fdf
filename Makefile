BINARY_NAME := fdf
VERSION := 1.0.0
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
