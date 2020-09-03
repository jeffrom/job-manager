SHELL := /bin/bash
TMPDIR := $(if $(TMPDIR),$(TMPDIR),"/tmp/")
GOPATH := $(shell go env GOPATH)

ctl_bin := jobctl
server_bin := job-manager
gofiles := $(wildcard *.go **/*.go **/**/*.go **/**/**/*.go)
protofiles := $(wildcard proto/*.proto proto/**/*.proto proto/**/**/*.proto proto/**/**/**/*.proto)
prototargets := $(wildcard doc/doc.json *.pb.go **/*.pb.go **/**/*.pb.go **/**/**/*.pb.go)

write_jsonschema_bin := script/write_jsonschema.sh
args_schema_target := pkg/schema/args_schema.go
args_schema_deps := jsonschema/Args.json
data_schema_target := pkg/schema/data_schema.go
data_schema_deps := jsonschema/Data.json
result_schema_target := pkg/schema/result_schema.go
result_schema_deps := jsonschema/Result.json

buf := $(shell which buf)
protoc := $(shell which protoc)
protoc_gen_go = $(GOPATH)/bin/protoc-gen-go
protoc_gen_doc = $(GOPATH)/bin/protoc-gen-doc
gocoverutil := $(GOPATH)/bin/gocoverutil
staticcheck := $(GOPATH)/bin/staticcheck
gomodoutdated := $(GOPATH)/bin/go-mod-outdated
tulpa := $(GOPATH)/bin/tulpa

ifeq ($(buf),)
	buf = must-rebuild
endif
ifeq ($(protoc),)
	protoc = must-rebuild
endif

all: build

build: $(gen) $(gofiles)
	GO111MODULE=on go install ./...

.make/$(server_bin): .make $(gen) $(gofiles)
	GO111MODULE=on go build -o .make/$(server_bin) ./cmd/$(server_bin)

.make/$(ctl_bin): .make $(gen) $(gofiles)
	GO111MODULE=on go build -o .make/$(ctl_bin) ./cmd/$(ctl_bin)

.make:
	mkdir -p .make

.PHONY: clean
clean:
	git clean -x -f

.PHONY: test
test: $(gen) $(gofiles) | $(staticcheck) $(buf)
	GO111MODULE=on go test -short ./...

.PHONY: test.lint
test.lint: $(gen) $(gofiles) | $(staticcheck) $(buf)
	GO111MODULE=on $(staticcheck) -f stylish -checks all ./...
	$(buf) check lint

.PHONY: test.cover
test.cover: $(gen) $(gofiles) | $(gocoverutil)
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

gen: gen.proto gen.jsonschema

gen.proto: $(prototargets)

$(prototargets): $(protofiles) | $(protoc_gen_go) $(protoc_gen_doc)
	protoc -I=proto --go_out=${GOPATH}/src ${protofiles}
	protoc -I=proto --doc_opt=json,doc.json --doc_out=doc ${protofiles}

gen.jsonschema: $(args_schema_target) $(data_schema_target) $(result_schema_target)

$(args_schema_target): $(write_jsonschema_bin) $(args_schema_deps)
	script/write_jsonschema.sh Args argSchemaRaw $(args_schema_target)

$(data_schema_target): $(write_jsonschema_bin) $(data_schema_deps)
	script/write_jsonschema.sh Data dataSchemaRaw $(data_schema_target)

$(result_schema_target): $(write_jsonschema_bin) $(result_schema_deps)
	script/write_jsonschema.sh Result resultSchemaRaw $(result_schema_target)

.PHONY: dev
dev:
	tulpa --ignore proto --ignore .make --ignore doc "make .make/$(server_bin) && .make/$(server_bin)"

$(gocoverutil):
	GO111MODULE=off go get github.com/AlekSi/gocoverutil

$(staticcheck):
	cd $(TMPDIR) && GO111MODULE=on go get honnef.co/go/tools/cmd/staticcheck@2020.1.5

$(gomodoutdated):
	GO111MODULE=off go get github.com/psampaz/go-mod-outdated

$(protoc_gen_go):
	GO111MODULE=off go get -u google.golang.org/protobuf/cmd/protoc-gen-go

$(protoc_gen_doc):
	GO111MODULE=off go get -u github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc

$(buf):
	@echo "Please install buf: https://buf.build/docs/installation/"
	@exit 1

$(protoc):
	@echo "Please install protoc"
	@exit 1

