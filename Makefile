APPNAME = abelwatch
VERSION=0.0.1
TESTFLAGS=-v -cover -covermode=atomic -bench=.
TEST_COVERAGE_THRESHOLD=8.0

build:
	go build -tags netgo -ldflags "-w" -o ${APPNAME} .

build-linux:
	GOOS=linux GOARCH=amd64 go build -tags netgo -ldflags "-w -s -X main.APP_VERSION=${VERSION}" -v -o ${APPNAME}-linux-amd64 .

build-mac:
	GOOS=darwin GOARCH=amd64 go build -tags netgo -ldflags "-w -s -X main.APP_VERSION=${VERSION}" -v -o ${APPNAME}-darwin-amd64 .

build-all: build-mac build-linux

clean:
	rm -f ${APPNAME}
	rm -f ${APPNAME}-linux-amd64
	rm -f ${APPNAME}-darwin-amd64

all: setup
	build
	install

setup:
	go get -u github.com/wadey/gocovmerge
	glide install

test:
	go test ${TESTFLAGS} -coverprofile=abel.txt github.com/indix/abelwatch/abel
	go test ${TESTFLAGS} -coverprofile=main.txt github.com/indix/abelwatch/
