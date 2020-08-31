package web

import (
	"context"
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi"
	chimw "github.com/go-chi/chi/middleware"
	"github.com/jeffrom/job-manager/pkg/web/handler"
	"github.com/jeffrom/job-manager/pkg/web/middleware"
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

func NewControllerRouter(cfg middleware.Config) (chi.Router, error) {
	r := chi.NewRouter()

	debugRoutes(r)

	logger := middleware.NewLogger(cfg.Logger)
	logger.Info().Interface("config", cfg).Msg("new router")

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

		r.Route("/internal", func(r chi.Router) {
			r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			})
			r.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			})
		})

		r.Route("/api/v1", func(r chi.Router) {
			r.Use(middleware.Backend(cfg.GetBackend()))

			r.Route("/jobs", func(r chi.Router) {
				r.Get("/", handler.Func(handler.ListQueues))

				r.Route("/{queueName}", func(r chi.Router) {
					r.Put("/", handler.Func(handler.SaveQueue))
					r.Delete("/", handler.Func(handler.DeleteQueue))

					r.Post("/enqueue", handler.Func(handler.EnqueueJobs))
					r.Post("/dequeue", handler.Func(handler.DequeueJobs))
				})

				r.Post("/ack", handler.Func(handler.Ack))
				r.Post("/dequeue", handler.Func(handler.DequeueJobs))
			})

			r.Route("/job/{jobID}", func(r chi.Router) {
				r.Get("/", handler.Func(handler.GetJobByID))
			})
		})
	})

	return r, nil
}
