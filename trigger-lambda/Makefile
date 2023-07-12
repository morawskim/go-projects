.PHONY: run
run:
	cd frontend && npm run build
	. .env && export LAMBDA_RUNTIME_URL && go run ./main.go
