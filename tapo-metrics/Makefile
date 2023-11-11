.PHONY: run publish-image build-image

run:
	. .env && export TAPO_IP TAPO_USERNAME TAPO_PASSWORD && go run ./

build-image:
	docker build -t ghcr.io/morawskim/tapo-metrics:latest .

publish-image: build-image
	teller run --  bash -c 'env | grep GHCR_TOKEN | cut -d'=' -f2 | docker login ghcr.io -u morawskim --password-stdin'
	docker push ghcr.io/morawskim/tapo-metrics:latest
