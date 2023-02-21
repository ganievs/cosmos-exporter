default: build

build:
	go build cmd/cosmos-exporter.go

install:
	go install cmd/cosmos-exporter.go

lint:
	golangci-lint run --fix ./...
