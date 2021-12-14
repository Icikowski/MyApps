.PHONY: clean build build-static package
.SILENT: ${.PHONY}

all: build-static package

get:
	(cd application; go mod download -x)

clean:
	rm -f packaging/myapps application/myapps
	rm -f myapps.tgz packaging/myapps.tgz
	(cd application; go clean -cache)

GIT_COMMIT := $(shell git rev-parse --short HEAD)
GIT_TAG := $(shell echo $${CURRENT_TAG:-$$(git describe --abbrev=0 | sed "s/v//")})
BASE_FLAGS := -X 'icikowski.pl/myapps/cli.version=${GIT_TAG}' -X 'icikowski.pl/myapps/cli.gitCommit=${GIT_COMMIT}'

build: clean
	(cd application; go build -ldflags "${BASE_FLAGS}" .)

build-static: clean
	(cd application; env CGO_ENABLED=0 go build -ldflags "${BASE_FLAGS} -w -extldflags '-static'" .)

package: build-static
	cp application/myapps packaging/myapps
	(cd packaging; tar -czf myapps.tgz myapps install.sh)
	mv packaging/myapps.tgz .