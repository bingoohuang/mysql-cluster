.PHONY: default test
all: default test

gosec:
	go get github.com/securego/gosec/cmd/gosec

sec:
	@gosec ./...
	@echo "[OK] Go security check was completed!"

proxy:
	export GOPROXY=https://goproxy.cn

default: proxy
	go fmt ./...&&revive .&&goimports -w .&&golangci-lint run --enable-all&&go install -ldflags="-s -w" ./...

install: proxy
	go install -ldflags="-s -w" ./...

test: proxy
	go test ./...


linux: proxy
	env GOOS=linux GOARCH=amd64 go install -ldflags="-s -w" ./...
	upx ~/go/bin/linux_amd64/mci