// Package bepg implements backend.Interface using Postgres.
package bepg

import (
	"context"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	db  *sqlx.DB
	cfg Config
	// TODO needs to be able to get the logger from ctx
}

type ProviderFunc func(p *Postgres) *Postgres

// example dsn:
// user=jack password=secret host=pg.example.com port=5432 dbname=mydb sslmode=verify-ca
func New(providers ...ProviderFunc) *Postgres {
	return &Postgres{}
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

	db, err := sqlx.ConnectContext(ctx, "pgx", pg.cfg.DSN())
	if err != nil {
		return err
	}
	pg.db = db
	return nil
}

func (pg *Postgres) getConn(ctx context.Context) (sqlx.ExtContext, error) {
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

func Reset(ctx context.Context) error {
	return nil
}

func (pg *Postgres) GetSetJobKeys(ctx context.Context, keys []string) (bool, error) {
	return false, nil
}

func (pg *Postgres) DeleteJobKeys(ctx context.Context, keys []string) error {
	return nil
}
