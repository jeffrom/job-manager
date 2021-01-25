package web

import (
	"context"
	"sync"

	"github.com/jeffrom/job-manager/pkg/backend"
	"github.com/jeffrom/job-manager/pkg/internal"
	"github.com/jeffrom/job-manager/pkg/web/middleware"
)

type Processor struct {
	cfg         middleware.Config
	be          backend.Interface
	reaper      backend.Reaper
	invalidator backend.Invalidator
}

func NewProcessor(cfg middleware.Config, be backend.Interface) *Processor {
	var reaper backend.Reaper
	if r, ok := be.(backend.Reaper); ok {
		reaper = r
	} else {
		reaper = &backend.NoopReaper{}
	}

	var invalidator backend.Invalidator
	if iv, ok := be.(backend.Invalidator); ok {
		invalidator = iv
	} else {
		invalidator = &backend.NoopInvalidator{}
	}

	return &Processor{
		cfg:         cfg,
		be:          be,
		reaper:      reaper,
		invalidator: invalidator,
	}
}

// Run starts a reaper and an invalidator goroutine, and returns any error that
// occurs while stopping.
func (p *Processor) Run(ctx context.Context) error {
	ctx = internal.EnsureTimeProvider(ctx)
	log := p.cfg.Logger
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info().Msg("start invalidator")
		for {
			if checkDone(ctx.Done()) {
				log.Info().Msg("invalidator shut down")
				return
			}
			internal.IgnoreError(p.RunInvalidateOnce(ctx))
			if checkDone(ctx.Done()) {
				log.Info().Msg("invalidator shut down")
				return
			}
			internal.IgnoreError(internal.Sleep(ctx, p.cfg.InvalidateInterval))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info().Msg("start reaper")
		for {
			if checkDone(ctx.Done()) {
				log.Info().Msg("reaper shut down")
				return
			}
			internal.IgnoreError(p.RunReaperOnce(ctx))
			if checkDone(ctx.Done()) {
				log.Info().Msg("reaper shut down")
				return
			}
			internal.IgnoreError(internal.Sleep(ctx, p.cfg.ReapInterval))
		}
	}()

	wg.Wait()
	return nil
}

func (p *Processor) RunOnce(ctx context.Context) error {
	ctx = internal.EnsureTimeProvider(ctx)
	if err := p.RunInvalidateOnce(ctx); err != nil {
		return err
	}
	if err := p.RunReaperOnce(ctx); err != nil {
		return err
	}
	return nil
}

func (p *Processor) RunInvalidateOnce(ctx context.Context) error {
	ctx = internal.EnsureTimeProvider(ctx)
	log := p.cfg.Logger
	log.Debug().Msg("start invalidation")
	err := p.invalidator.InvalidateJobs(ctx)
	if err != nil {
		log.Error().Err(err).Msg("invalidation failed")
	} else {
		log.Debug().Msg("invalidation complete")
	}
	return err
}

func (p *Processor) RunReaperOnce(ctx context.Context) error {
	ctx = internal.EnsureTimeProvider(ctx)
	log := p.cfg.Logger
	log.Debug().Msg("start reap")
	err := p.reaper.Reap(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("reap failed")
	} else {
		log.Debug().Msg("reap complete")
	}
	return err
}

func checkDone(done <-chan struct{}) bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}
