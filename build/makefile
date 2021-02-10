version = $(shell  cat ./../VERSION)

release:
	BUILD_VERSION=$(version)
	cd .. && goreleaser --skip-publish --rm-dist
