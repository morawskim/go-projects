.PHONY: run build publish frontend

frotnend:
	cd frontend && PUBLIC_URL='{{.BasePath}}' npm run build

run: frotnend
	. .env && export LAMBDA_RUNTIME_URL && go run ./

build: frotnend
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o trigger-lambda-amd64 ./
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o trigger-lambda-arm64 ./

publish: build
	docker buildx build --platform linux/amd64,linux/arm64 -t ghcr.io/morawskim/trigger-lambda:latest -t morawskim/trigger-lambda:latest . --push
