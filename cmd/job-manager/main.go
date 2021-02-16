package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"

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

	proc := web.NewProcessor(*cfg, be)
	var srv *http.Server
	var ln net.Listener
	if len(os.Args) == 1 {
		srv, ln, err = createSrv(cfg, be)
		if err != nil {
			panic(err)
		}
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	ctx, done := context.WithCancel(context.Background())
	defer done()
	go func() {
		<-sigs
		if srv != nil {
			srv.Shutdown(ctx)
		}
		done()
	}()

	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "reaper":
			if err := proc.Run(ctx); err != nil {
				panic(err)
			}
			return
		case "reap":
			if err := proc.RunOnce(ctx); err != nil {
				panic(err)
			}
			return
		}
		panic("unknown arg: " + os.Args[1])
	}

	wg := sync.WaitGroup{}
	if dev := os.Getenv("REAPER"); dev != "" {
		cfg.Logger.Info().Msg("running local reaper/invalidator")
		wg.Add(1)
		go func() {
			if err := proc.Run(ctx); err != nil {
				panic(err)
			}
			wg.Done()
		}()
	}

	wg.Add(1)
	go func() {
		if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
		wg.Done()
	}()

	wg.Wait()
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

func createSrv(cfg *middleware.Config, be backend.Interface) (*http.Server, net.Listener, error) {
	h, err := web.NewControllerRouter(*cfg, be)
	if err != nil {
		return nil, nil, err
	}

	addr := ":1874"
	if cfg.Host != "" {
		addr = cfg.Host
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: h,
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, err
	}

	return srv, ln, nil
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
