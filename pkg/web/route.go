package web

import (
	"context"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/go-chi/chi"
	chimw "github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/jeffrom/job-manager/pkg/backend"
	"github.com/jeffrom/job-manager/pkg/internal"
	"github.com/jeffrom/job-manager/pkg/web/handler"
	"github.com/jeffrom/job-manager/pkg/web/middleware"
	"github.com/jeffrom/job-manager/release"
)

func debugRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.HandleFunc("/debug/pprof", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, r.RequestURI+"/", http.StatusMovedPermanently)
		})
		r.HandleFunc("/debug/pprof/*", pprof.Index)
		r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		r.HandleFunc("/debug/pprof/profile", pprof.Profile)
		r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		r.HandleFunc("/debug/pprof/trace", pprof.Trace)
	})
}

func NewControllerRouter(cfg middleware.Config, be backend.Interface) (chi.Router, error) {
	r := chi.NewRouter()

	debugRoutes(r)

	logger := cfg.Logger
	logger.Info().
		Str("v", release.Version).
		Str("commit", release.Commit).
		Interface("config", cfg).
		Msg("new router")

	r.Group(func(r chi.Router) {
		r.Use(
			// set config, don't want to deal w/ making a new package for this
			func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					ctx := context.WithValue(r.Context(), middleware.ConfigKey, cfg)
					next.ServeHTTP(w, r.WithContext(ctx))
				})
			},
			chimw.RealIP,
			chimw.RequestID,
			chimw.StripSlashes,
			logger.Middleware,
			chimw.Recoverer,
		)
		r.MethodNotAllowed(handler.Func(handler.MethodNotAllowed))
		r.NotFound(handler.Func(handler.NotFound))

		r.Handle("/metrics", promhttp.Handler())

		r.Route("/internal", func(r chi.Router) {
			r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			})
			r.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			})
		})

		r.Route("/api/v1", func(r chi.Router) {
			r.Use(
				middleware.Time(internal.Time(0), internal.NewTicker(1*time.Second)),
				backend.Middleware(be),
				// middleware.Backend(cfg.GetBackend()),
			)

			if hp, ok := be.(backend.HandlerProvider); ok {
				r.Route("/backend", func(r chi.Router) {
					r.Handle("/*", hp.Handler())
				})
			}

			r.Route("/", func(r chi.Router) {
				if mwp, ok := be.(backend.MiddlewareProvider); ok {
					r.Use(mwp.Middleware())
				}

				r.Get("/stats/{queueName}", handler.Func(handler.Stats))
				r.Get("/stats", handler.Func(handler.Stats))

				r.Route("/queues", func(r chi.Router) {
					r.Get("/", handler.Func(handler.ListQueues))

					r.Route("/{queueID}", func(r chi.Router) {
						r.Get("/", handler.Func(handler.GetQueueByID))
						r.Put("/", handler.Func(handler.SaveQueue))
						r.Delete("/", handler.Func(handler.DeleteQueue))

						enqueueHandler := &handler.EnqueueJobs{}
						r.Post("/enqueue", enqueueHandler.ServeHTTP)
						r.Post("/dequeue", handler.Func(handler.DequeueJobs))
					})
				})

				r.Route("/jobs", func(r chi.Router) {
					r.Post("/dequeue", handler.Func(handler.DequeueJobs))
					r.Post("/ack", handler.Func(handler.Ack))

					r.Route("/{jobID}", func(r chi.Router) {
						r.Get("/", handler.Func(handler.GetJobByID))
						r.Get("/queue", handler.Func(handler.GetQueueByJobID))
					})
					r.Get("/", handler.Func(handler.ListJobs))
				})
			})
		})
	})

	return r, nil
}
