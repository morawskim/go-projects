.PHONY: build run run-in-k8s

build:
	docker build -t admission-webhook-demo ./
	minikube image rm admission-webhook-demo
	minikube image load admission-webhook-demo

run:
	go run ./ -tlsKey certs/tls.key -tlsCert certs/tls.crt

run-in-k8s:
	go run ./
