.PHONY: default test
all: default test

gosec:
	go get github.com/securego/gosec/cmd/gosec
sec:
	@gosec ./...
	@echo "[OK] Go security check was completed!"

init:
	export GOPROXY=https://goproxy.cn

default: init
	gofmt -s -w .&&go mod tidy&&go fmt ./...&&revive .&&goimports -w .&&golangci-lint run --enable-all&&go install -ldflags="-s -w" ./...

install: init packr2
	go install -ldflags="-s -w" ./...

test: init
	go test ./...

linux: init packr2
	GOOS=linux GOARCH=amd64 go install -ldflags="-s -w" ./...
	# upx ~/go/bin/linux_amd64/otterbeat

packr2:
	cd cmd/otterbeat ; go install github.com/gobuffalo/packr/v2; packr2