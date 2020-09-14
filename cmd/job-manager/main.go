package main

import (
	"net"
	"net/http"

	"github.com/jeffrom/job-manager/pkg/backend/bememory"
	"github.com/jeffrom/job-manager/pkg/web"
	"github.com/jeffrom/job-manager/pkg/web/middleware"
)

func main() {
	h, err := web.NewControllerRouter(middleware.NewConfig(), bememory.New())
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
