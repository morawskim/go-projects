.PHONY: build-image publish-image

build-image:
	docker build -t ghcr.io/morawskim/price-tracker:latest .

publish-image: build-image
	teller run --  bash -c 'env | grep GHCR_TOKEN | cut -d'=' -f2 | docker login ghcr.io -u morawskim --password-stdin'
	docker push ghcr.io/morawskim/price-tracker:latest
