# This is how we want to name the binary output
BINARY=bin/graphss

# These are the values we want to pass for version and build
# git tag 1.0.0
# git commit -am "One more commit after the tags"
# VERSION=`git describe --tags`
VERSION=v1.0.0
BUILD=`git rev-parse HEAD`

# Setup ldflags option for go build
LDFLAGS=-ldflags="-X 'main.version=${VERSION}' -X 'main.commit=${BUILD}'"

run: gen
	@echo "running"
	go run ${LDFLAGS} pkg/cmd/graphss/main.go --config ./conf/config.toml

build: gen
	@echo "building a binary"
	go build ${LDFLAGS} -o ${BINARY} pkg/cmd/graphss/main.go

gen:
	@echo "generating wire files"
	wire ./pkg/http

clean:
	if [ -f ${BINARY} ]; then rm ${BINARY} ; fi

.PHONY: clean install