.PHONY: release

release:
	teller run --  bash -c 'env | grep GHCR_TOKEN | cut -d'=' -f2 | docker login ghcr.io -u morawskim --password-stdin'
	~/go/bin/goreleaser release --clean
