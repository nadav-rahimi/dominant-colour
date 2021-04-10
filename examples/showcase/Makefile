BINARY_NAME=showcase

build:
		go mod download
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)_linux
		CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME)_windows
clean:
		go clean
		rm -f bin/$(BINARY_NAME)_linux
		rm -f bin/$(BINARY_NAME)_windows
bench:
		go test -bench=. -benchtime=10s