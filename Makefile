.PHONY: all build test lint clean coverage

all: build

build:
	go build -o hyperagent main.go

test:
	go test -v ./...

lint:
	/root/go/bin/golangci-lint run ./...

clean:
	rm -f hyperagent

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
