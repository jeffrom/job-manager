// +build tools

// Package tools is intended to version protoc-gen-go in go.mod.
package tools

import (
	_ "github.com/jteeuwen/go-bindata/go-bindata"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
