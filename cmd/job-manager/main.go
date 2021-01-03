package main

import (
	"errors"
	"net"
	"net/http"
	"os"

	"github.com/jeffrom/job-manager/pkg/backend"
	"github.com/jeffrom/job-manager/pkg/backend/bememory"
	"github.com/jeffrom/job-manager/pkg/backend/bepostgres"
	"github.com/jeffrom/job-manager/pkg/config"
	"github.com/jeffrom/job-manager/pkg/logger"
	"github.com/jeffrom/job-manager/pkg/web"
	"github.com/jeffrom/job-manager/pkg/web/middleware"
)

func main() {
	cfgFlags := &middleware.Config{}
	icfg, err := config.MergeEnvFlags(cfgFlags, &middleware.ConfigDefaults)
	if err != nil {
		panic(err)
	}

	cfg := icfg.(*middleware.Config)
	cfg.ResetLogOutput(os.Stdout)

	be, err := backendWithConfig(cfg.Backend, cfg.Logger)
	if err != nil {
		panic(err)
	}

	if err := run(cfg, be); err != nil {
		panic(err)
	}
}

func backendWithConfig(name string, logger *logger.Logger) (backend.Interface, error) {
	switch name {
	case "postgres":
		ibecfg, err := config.MergeEnvFlags(&bepostgres.Config{}, &bepostgres.DefaultConfig)
		if err != nil {
			return nil, err
		}
		becfg := ibecfg.(*bepostgres.Config)
		becfg.Logger = logger
		return bepostgres.New(bepostgres.WithConfig(*becfg)), nil
	case "memory":
		return bememory.New(), nil
	}

	return nil, errors.New("unknown backend: " + name)
}

func run(cfg *middleware.Config, be backend.Interface) error {
	h, err := web.NewControllerRouter(*cfg, be)
	if err != nil {
		return err
	}

	host := ":1874"
	if cfg.Host != "" {
		host = cfg.Host
	}
	ln, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}

	if err := http.Serve(ln, h); err != nil {
		return err
	}

	return nil
}
