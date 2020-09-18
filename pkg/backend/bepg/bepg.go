// Package bepg implements backend.Interface using Postgres.
package bepg

import (
	"context"

	// "github.com/jackc/pgx/v4/pgxpool"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zerologadapter"
	"github.com/jackc/pgx/v4/stdlib"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jeffrom/job-manager/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type sqlxer interface {
	sqlx.ExtContext
	sqlx.PreparerContext
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
}

type Postgres struct {
	db  *sqlx.DB
	cfg Config
	// TODO needs to be able to get the logger from ctx
}

type ProviderFunc func(p *Postgres) *Postgres

// example dsn:
// user=jack password=secret host=pg.example.com port=5432 dbname=mydb sslmode=verify-ca
func New(providers ...ProviderFunc) *Postgres {
	pg := &Postgres{}
	for _, provider := range providers {
		pg = provider(pg)
	}
	return pg
}

func WithConfig(cfg Config) ProviderFunc {
	return func(p *Postgres) *Postgres {
		p.cfg = cfg
		return p
	}
}

func (pg *Postgres) ensureConn(ctx context.Context) error {
	if pg.db != nil {
		return nil
	}
	reqlog := logger.RequestLogFromContext(ctx)
	reqlog.Bool("dbconnect", true)

	log := logger.FromContext(ctx)

	// pgCfg := &pgx.ConnConfig{
	// 	Config: pgconn.Config{
	// 		Database: pg.cfg.Database,
	// 		Host:     pg.cfg.Host,
	// 		Port:     uint16(pg.cfg.Port),
	// 		User:     pg.cfg.User,
	// 		Password: pg.cfg.Password,
	// 	},
	// 	Logger: zerologadapter.NewLogger(pg.cfg.Logger.Logger),
	// }

	dsn := pg.cfg.DSN()
	log.Debug().Str("dsn", dsn).Msg("connecting to postgres")
	pgCfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		return err
	}
	pgCfg.Logger = zerologadapter.NewLogger(pg.cfg.Logger.Logger)
	if pg.cfg.Debug {
		pgCfg.LogLevel = pgx.LogLevelTrace
	}

	connStr := stdlib.RegisterConnConfig(pgCfg)

	db, err := sqlx.ConnectContext(ctx, "pgx", connStr)
	if err != nil {
		return err
	}
	pg.db = db
	return nil
}

func (pg *Postgres) getConn(ctx context.Context) (sqlxer, error) {
	if tx := getTx(ctx); tx != nil {
		return tx, nil
	}
	if err := pg.ensureConn(ctx); err != nil {
		return nil, nil
	}
	return pg.db, nil
}

func (pg *Postgres) Ping(ctx context.Context) error {
	if err := pg.ensureConn(ctx); err != nil {
		return err
	}
	return nil
}

func (pg *Postgres) Reset(ctx context.Context) error {
	return nil
}

func (pg *Postgres) GetSetJobKeys(ctx context.Context, keys []string) (bool, error) {
	return false, nil
}

func (pg *Postgres) DeleteJobKeys(ctx context.Context, keys []string) error {
	return nil
}
