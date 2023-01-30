.PHONY: build
build:
	go build cmd/tg_bot_api/main.go

.PHONY: test
	go test -v ./...