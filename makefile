all: build
# This how we want to name the binary output
BINARY=httpgo

# These are the values we want to pass for VERSION and BUILD
# git tag 1.0.1
# git commit -am "One more change after the tags"
VERSION=`git describe --tags`
BUILD=`date +%FT%T%z`
BRANCH=`git branch --show-current`

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-s -w -X main.Version=${VERSION} -X main.Build=${BUILD} -X main.Branch=${BRANCH}"
mod:
	go mod tidy
#	chown ruslan:progs go.*
# Builds the project
run:
	go run -race ${LDFLAGS} httpgo.go
#httpgo only
httpgo_all:
	go mod tidy
	chown ruslan:progs go.*
	go generate
	go build -i ${LDFLAGS} -o ${BINARY}
	systemctl restart httpgo
	sleep 20
	journalctl -u httpgo --since="30 second ago" -o cat
# Builds the project
build:
	go build ${LDFLAGS} -o ${BINARY}
# Builds dev server
dev:
	go build -race ${LDFLAGS} -o ${BINARY}_dev
# test
test:
	go test -v ./... > last_test.log
# Builds the project
linux:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -i ${LDFLAGS} -o ${BINARY}

# Cleans our project: deletes binaries

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: clean install