package main

import (
	"net"
	"net/http"
	"os"

	"github.com/jeffrom/job-manager/pkg/backend/bepostgres"
	"github.com/jeffrom/job-manager/pkg/logger"
	"github.com/jeffrom/job-manager/pkg/web"
	"github.com/jeffrom/job-manager/pkg/web/middleware"
)

func main() {
	log := logger.New(os.Stdout, false)

	// be := beredis.New()
	// be := bememory.New()
	becfg := bepostgres.DefaultConfig
	becfg.Logger = log
	be := bepostgres.New(bepostgres.WithConfig(becfg))

	cfg := middleware.NewConfig()
	cfg.Logger = log

	h, err := web.NewControllerRouter(cfg, be)
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
