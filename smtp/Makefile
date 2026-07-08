.PHONY: package build

build:
	go build -o gosmtp ./server.go

package: build
	cd pkg/systemd/ && jsonnet service.libsonnet -S -m .
	docker run --rm -v $(PWD):/tmp -e PACKAGER=deb -w /tmp goreleaser/nfpm pkg --packager deb --target /tmp
	docker run --rm -v $(PWD):/tmp -e PACKAGER=rpm -w /tmp goreleaser/nfpm pkg --packager rpm --target /tmp
