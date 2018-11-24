GO=GO15VENDOREXPERIMENT="1" CGO_ENABLED=0 GO111MODULE=on go
GO_CGO=GO15VENDOREXPERIMENT="1" CGO_ENABLED=1 GO111MODULE=on go
GOTEST=GO15VENDOREXPERIMENT="1" CGO_ENABLED=1 GO111MODULE=on go test # go race detector requires cgo

PACKAGES := $$(go list ./...| grep -vE 'vendor|agent' )

FILES     := $$(find . -name "*.go" | grep -vE "vendor")
GOFILTER := grep -vE 'vendor|render.Delims|bindata_assetfs|testutil|\.pb\.go'
GOCHECKER := $(GOFILTER) | awk '{ print } END { if (NR > 0) { exit 1 } }'
GOLINT := go list ./... | grep -vE 'vendor' | xargs -L1 -I {} golint {} 2>&1 | $(GOCHECKER)

LDFLAGS += -X "github.com/pingcap/schrodinger/util.BuildTS=$(shell date -u '+%Y-%m-%d %I:%M:%S')"
LDFLAGS += -X "github.com/pingcap/schrodinger/util.BuildHash=$(shell git rev-parse HEAD)"

GOBUILD=$(GO) build -ldflags '$(LDFLAGS)'
GOBUILD_CGO=$(GO_CGO) build -ldflags '$(LDFLAGS)'

GOMOD := -mod=vendor

default: all

all: server client

server:
        $(GOBUILD_CGO) $(GOMOD) -o bin/server server/*.go
client:
		$(GOBUILD_CGO) $(GOMOD) -o bin/client client/*.go

check:
        GO111MODULE=on go get github.com/golang/lint/golint

        @echo "vet"

        @go vet -all -shadow ./...

        @echo "golint"
        @ $(GOLINT)
        @echo "gofmt"
        @gofmt -s -l -w $(FILES) 2>&1 | $(GOCHECKER)

clean:
        @rm -rf bin/*

.PHONY: all client server