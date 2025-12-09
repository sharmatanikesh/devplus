.PHONY: run tidy build-mac

run:
	go run cmd/server/main.go

tidy:
	go mod tidy

build-mac:
	GOOS=darwin GOARCH=arm64 go build -o bin/devplus-backend cmd/server/main.go
