.PHONY: tests
tests:
	docker run --rm -v $(PWD):/app -w /app --entrypoint=promtool prom/prometheus:v3.2.1 test rules alerts_test.yml
