.PHONY: package build

build:
	go build -o gosmtp ./server.go

package:
# 	docker run --rm -v $(PWD):/tmp -w /tmp goreleaser/nfpm pkg --packager rpm --target /tmp
	docker run --rm -v $(PWD):/tmp -w /tmp goreleaser/nfpm pkg --packager deb --target /tmp
