.PHONY: build-wasm build-server run

build-wasm:
	GOOS=js GOARCH=wasm go build -o ./assets/functions.wasm ./cmd/wasm
	cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" ./assets/

build-server:
	go build -o server ./cmd/server

run:
	go run ./cmd/server/main.go
