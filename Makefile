SHELL := /bin/bash
TMPDIR := $(if $(TMPDIR),$(TMPDIR),"/tmp/")
GOPATH := $(shell go env GOPATH)

bins := $(GOPATH)/bin/jobctl
gofiles := $(wildcard *.go **/*.go **/**/*.go **/**/**/*.go)
protofiles := $(wildcard proto/*.proto proto/**/*.proto proto/**/**/*.proto proto/**/**/**/*.proto)
prototargets := $(wildcard *.pb.go **/*.pb.go **/**/*.pb.go **/**/**/*.pb.go)

buf := $(shell which buf)
protoc_gen_go = $(GOPATH)/bin/protoc-gen-go
gocoverutil := $(GOPATH)/bin/gocoverutil
staticcheck := $(GOPATH)/bin/staticcheck
gomodoutdated := $(GOPATH)/bin/go-mod-outdated

ifeq ($(buf),)
	buf = must-rebuild
endif

all: build

build: $(bins)

$(bins): $(gen) $(gofiles)
	GO111MODULE=on go install ./cmd/...

.PHONY: clean
clean:
	git clean -x -f

.PHONY: test
test: $(gen)
	GO111MODULE=on go test -cover -race ./...

.PHONY: test.lint
test.lint: $(gen) | $(staticcheck)
	GO111MODULE=on $(staticcheck) -f stylish -checks all ./...

.PHONY: test.cover
test.cover: $(gen) | $(gocoverutil)
	$(gocoverutil) -coverprofile=cov.out test -covermode=count ./... \
		2> >(grep -v "no packages being tested depend on matches for pattern" 1>&2) \
		| sed -e 's/of statements in .*/of statements/'
	@echo -n "total: "; go tool cover -func=cov.out | tail -n 1 | sed -e 's/\((statements)\|total:\)//g' | tr -s "[:space:]"

.PHONY: test.outdated
test.outdated: $(gomodoutdated)
	GO111MODULE=on go list -u -m -json all | go-mod-outdated -direct

.PHONY: release.dryrun
release.dryrun:
	goreleaser --snapshot --skip-publish --rm-dist

.PHONY: release
release:
	goreleaser --rm-dist

gen: gen.proto

gen.proto: $(prototargets)

$(prototargets): $(protofiles) | $(protoc_gen_go) $(buf)
	buf protoc -I=proto --go_out=${GOPATH}/src ${protofiles}

$(gocoverutil):
	GO111MODULE=off go get github.com/AlekSi/gocoverutil

$(staticcheck):
	cd $(TMPDIR) && GO111MODULE=on go get honnef.co/go/tools/cmd/staticcheck@2020.1.5

$(gomodoutdated):
	GO111MODULE=off go get github.com/psampaz/go-mod-outdated

$(protoc_gen_go):
	GO111MODULE=off go get -u github.com/golang/protobuf/protoc-gen-go

$(buf):
	@echo "Please install buf: https://buf.build/docs/installation/"
	@exit 1

