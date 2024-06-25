.PHONY: run build-image publish-image
TAG := v0.1

run:
	. .env && export FILE_URL && go run main.go

build-image:
	docker build -t ghcr.io/morawskim/webcrypto:$(TAG) .

publish-image: build-image
	teller run --  bash -c 'env | grep GHCR_TOKEN | cut -d'=' -f2 | docker login ghcr.io -u morawskim --password-stdin'
	docker push ghcr.io/morawskim/webcrypto:$(TAG)
