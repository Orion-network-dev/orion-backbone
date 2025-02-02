build-apk:
	goreleaser release --snapshot --clean
.PHONY: build-apk