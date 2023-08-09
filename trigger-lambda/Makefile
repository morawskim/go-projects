.PHONY: run build publish
run:
	cd frontend && npm run build
	. .env && export LAMBDA_RUNTIME_URL && go run ./main.go

build:
	cd frontend && npm run build
	GOOS=linux GOARCH=amd64 go build -o trigger-lambda-amd64 ./main.go
	GOOS=linux GOARCH=arm64 go build -o trigger-lambda-arm64 ./main.go

publish: build
	docker buildx build --platform linux/amd64,linux/arm64 -t ghcr.io/morawskim/trigger-lambda:latest -t morawskim/trigger-lambda:latest . --push
