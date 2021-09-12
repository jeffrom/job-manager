FROM golang:1.17.1 as builder

WORKDIR /build

COPY go.mod go.sum /build/
COPY mjob/go.mod mjob/go.sum /build/mjob/
RUN set -x; go mod download && cd mjob && go mod download

COPY . /build

ARG VERSION=next
ARG COMMIT=none
ARG DATE=none
RUN set -x; CGO_ENABLED=0 go build -o job-manager.bin -ldflags "-s -w -X github.com/jeffrom/job-manager/release.Version=${VERSION} -X github.com/jeffrom/job-manager/release.Commit=${COMMIT} -X github.com/jeffrom/job-manager/release.Date=${DATE}" ./cmd/job-manager

FROM scratch

COPY --from=builder /build/job-manager.bin /usr/local/bin/job-manager

ENTRYPOINT ["job-manager"]
