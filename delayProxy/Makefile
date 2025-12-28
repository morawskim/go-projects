.PHONY: build-image publish-image

TAG := v0.1.0

build-image:
	docker build --pull -t ghcr.io/morawskim/delay-proxy:$(TAG) .

publish-image: build-image
	#teller run --  bash -c 'env | grep GHCR_TOKEN | cut -d'=' -f2 | docker login ghcr.io -u morawskim --password-stdin'
	docker push ghcr.io/morawskim/delay-proxy:$(TAG)
