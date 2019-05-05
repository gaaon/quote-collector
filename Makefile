test:
	go test -race -v $(shell go list ./pkg/... | grep -v /vendor/) -coverprofile=c.out && go tool cover -html=c.out

BINARY        ?= azr-manager
SOURCES        = $(shell find . -name '*.go')
IMAGE         ?= registry.opensource.zalan.do/teapot/$(BINARY)
VERSION       ?= $(shell git describe --tags --always --dirty)
BUILD_FLAGS   ?= -v
LDFLAGS       ?= -X github.com/kubernetes-incubator/external-dns/pkg/apis/externaldns.Version=$(VERSION) -w -s

build:
	CGO_ENABLED=0 go build -o bin/${BINARY} ${BUILD_FLAGS} cmd/quote-collector/*.go

collect:
	go run cmd/collect-people/main.go

run:
	go run cmd/quote-collector/main.go ${FILE_VERSION}
