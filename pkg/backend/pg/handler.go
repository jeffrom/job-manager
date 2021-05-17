package pg

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	bindata "github.com/golang-migrate/migrate/source/go_bindata"

	"github.com/jeffrom/job-manager/pkg/backend"
	"github.com/jeffrom/job-manager/pkg/backend/pg/migrations"
	"github.com/jeffrom/job-manager/pkg/logger"
)

func (pg *Postgres) Handler() http.Handler {
	return http.HandlerFunc(handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.URL.EscapedPath(), "/migrate") {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	log := logger.FromContext(ctx)
	pg := backend.FromMiddleware(ctx).(*Postgres)
	cfgCopy := pg.cfg
	cfgCopy.Database = "postgres"

	log.Info().Msg("running migration")
	connStr, err := registerConnConfig(cfgCopy.DSN(), pg.cfg.Logger.Logger, pg.cfg.Debug)
	if err != nil {
		logWriteError(w, log, err)
		return
	}

	conn, err := sql.Open("pgx", connStr)
	if err != nil {
		logWriteError(w, log, err)
		return
	}
	defer conn.Close()

	if _, err := conn.ExecContext(ctx, "CREATE DATABASE "+pg.cfg.Database); err != nil {
		log.Info().Err(err).Msg("ignoring error")
	}
	if err := conn.Close(); err != nil {
		log.Info().Err(err).Msg("ignoring error")
	}

	// connection on the database we just ensured exists
	dbConnStr, err := registerConnConfig(pg.cfg.DSN(), pg.cfg.Logger.Logger, pg.cfg.Debug)
	if err != nil {
		logWriteError(w, log, err)
		return
	}

	dbConn, err := sql.Open("pgx", dbConnStr)
	if err != nil {
		logWriteError(w, log, err)
		return
	}

	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	if err != nil {
		logWriteError(w, log, err)
		return
	}

	migrationData := bindata.Resource(migrations.AssetNames(),
		func(name string) ([]byte, error) {
			return migrations.Asset(name)
		},
	)
	migrationDriver, err := bindata.WithInstance(migrationData)
	if err != nil {
		logWriteError(w, log, err)
		return
	}

	m, err := migrate.NewWithInstance("go-bindata", migrationDriver, "postgres", driver)
	if err != nil {
		logWriteError(w, log, err)
		return
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logWriteError(w, log, err)
		return
	}
}

func logWriteError(w http.ResponseWriter, log *logger.Logger, err error) {
	log.Error().Err(err).Msg("migration failed")
	w.WriteHeader(http.StatusInternalServerError)
}
