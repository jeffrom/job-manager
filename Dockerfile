FROM golang:1.15.6 as builder

RUN mkdir /build
WORKDIR /build

COPY go.mod /build/
COPY go.sum /build/
COPY mjob /build/
RUN go mod download

COPY . /build

ARG VERSION=next
ARG COMMIT=none
ARG DATE=none
RUN CGO_ENABLED=0 go build -o job-manager.bin -ldflags "-s -w -X github.com/jeffrom/job-manager/release.Version=${VERSION} -X github.com/jeffrom/job-manager/release.Commit=${COMMIT} -X github.com/jeffrom/job-manager/release.Date=${DATE}" ./cmd/job-manager

FROM alpine:3.12.3

COPY --from=builder /build/job-manager.bin /usr/local/bin/job-manager

ENTRYPOINT ["job-manager"]
