version = $(shell  cat ./../VERSION)

release:
	export BUILD_VERSION=$(version)
	cd .. && goreleaser --skip-publish --rm-dist
