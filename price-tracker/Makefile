.PHONY: build-image publish-image

TAG := v0.1

build-image:
	docker build -t ghcr.io/morawskim/price-tracker:$(TAG) .

publish-image: build-image
	teller run --  bash -c 'env | grep GHCR_TOKEN | cut -d'=' -f2 | docker login ghcr.io -u morawskim --password-stdin'
	docker push ghcr.io/morawskim/price-tracker:$(TAG)
