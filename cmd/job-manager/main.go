package main

import (
	"net"
	"net/http"

	"github.com/jeffrom/job-manager/pkg/backend"
	"github.com/jeffrom/job-manager/pkg/web"
)

func main() {
	be := backend.NewMemory()
	h, err := web.NewControllerRouter(be)
	if err != nil {
		panic(err)
	}

	ln, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic(err)
	}

	if err := http.Serve(ln, h); err != nil {
		panic(err)
	}
}
