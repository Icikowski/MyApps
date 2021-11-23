.PHONY: clean build build-static package
.SILENT: ${.PHONY}

all: build-static package

clean:
	rm -f packaging/myapps application/myapps
	(cd application; go clean -cache)

build: clean
	(cd application; go build .)

build-static: clean
	(cd application; env CGO_ENABLED=0 go build -ldflags "-w -extldflags '-static'" .)

package: build-static
	cp application/myapps packaging/myapps
	(cd packaging; tar -czf myapps.tgz myapps install.sh)
	mv packaging/myapps.tgz .