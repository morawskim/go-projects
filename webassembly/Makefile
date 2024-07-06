.PHONY: build-wasm build-server run

build-wasm:
	GOOS=js GOARCH=wasm go build -o ./assets/functions.wasm ./cmd/wasm
	cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" ./assets/

build-server:
	go build -o server ./cmd/server

run: build-wasm
	go run ./cmd/server/main.go

build-via-tiny-go:
	docker run --rm -v $(PWD):/app -w /app tinygo/tinygo:0.32.0 tinygo build -o ./assets/functions.wasm -target wasm --no-debug ./cmd/wasm
	docker run --rm -v $(PWD):/app -w /app tinygo/tinygo:0.32.0 /bin/bash -c "cp /usr/local/tinygo/targets/wasm_exec.js ./assets/"

run-tinygo: build-via-tiny-go
	go run ./cmd/server/main.go
