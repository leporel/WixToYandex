export  BUILD_VERSION = $(shell  cat ./../VERSION)

release:
	cd .. && goreleaser --skip-publish --rm-dist
