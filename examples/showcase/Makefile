 # Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=showcase

build:
		# Get dependencies
		$(GOMOD) download

		# Build for linux and windows
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o bin/$(BINARY_NAME)_linux
		CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o bin/$(BINARY_NAME)_windows
clean:
		$(GOCLEAN)
		rm -f bin/$(BINARY_NAME)_linux
		rm -f bin/$(BINARY_NAME)_windows
