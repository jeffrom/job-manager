// Package bepg implements backend.Interface using Postgres.
package bepg

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	// "github.com/jackc/pgx/v4/pgxpool"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zerologadapter"
	"github.com/jackc/pgx/v4/stdlib"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/jeffrom/job-manager/pkg/logger"
)

type sqlxer interface {
	sqlx.ExtContext
	sqlx.PreparerContext
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
}

type Postgres struct {
	db  *sqlx.DB
	cfg Config
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

func (pg *Postgres) getLogger(ctx context.Context) *logger.Logger {
	log := logger.FromContext(ctx)
	if log == nil {
		log = pg.cfg.Logger
	}
	return log
}

func (pg *Postgres) ensureConn(ctx context.Context) error {
	if pg.db != nil {
		return nil
	}
	reqlog := logger.RequestLogFromContext(ctx)
	reqlog.Bool("dbconnect", true)

	log := pg.getLogger(ctx)

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
	if err := pg.ensureConn(ctx); err != nil {
		return err
	}
	rows, err := pg.db.QueryxContext(ctx, "SELECT tableowner, tablename FROM pg_tables WHERE tableowner != 'postgres'")
	if err != nil {
		return err
	}

	var tables []string
	for rows.Next() {
		islice := []interface{}{}
		if err := rows.Scan(islice); err != nil {
			return err
		}
		tables = append(tables, islice[0].(string))
	}

	for _, table := range tables {
		if _, err := pg.db.ExecContext(ctx, "TRUNCATE TABLE "+table); err != nil {
			return err
		}
	}
	return nil
}

func (pg *Postgres) GetSetJobKeys(ctx context.Context, keys []string) (bool, error) {
	c, err := pg.getConn(ctx)
	if err != nil {
		return false, err
	}
	q := "SELECT 't'::boolean FROM job_uniqueness WHERE key IN (?)"
	args := stringsToBytea(keys)
	q, iargs, err := sqlx.In(q, args)
	if err != nil {
		return false, err
	}
	rows, err := c.QueryxContext(ctx, c.Rebind(q), iargs...)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}
	for rows.Next() {
		return true, nil
	}

	stmt, err := c.PrepareContext(ctx, "INSERT INTO job_uniqueness (key) VALUES ($1)")
	if err != nil {
		return false, err
	}
	for _, arg := range iargs {
		if _, err := stmt.ExecContext(ctx, arg); err != nil {
			return false, err
		}
	}
	return false, nil
}

func (pg *Postgres) DeleteJobKeys(ctx context.Context, keys []string) error {
	c, err := pg.getConn(ctx)
	if err != nil {
		return err
	}
	stmt, err := c.PrepareContext(ctx, "DELETE FROM job_uniqueness WHERE key = $1")
	if err != nil {
		return err
	}
	for _, key := range keys {
		if _, err := stmt.ExecContext(ctx, key); err != nil {
			return err
		}
	}
	return nil
}

func sqlFields(fields ...string) (string, string) {
	cols := strings.Join(fields, ", ")
	args := make([]string, len(fields))
	for i, field := range fields {
		args[i] = ":" + field
	}
	return cols, strings.Join(args, ", ")
}

func stringsToBytea(vals []string) [][]byte {
	res := make([][]byte, len(vals))
	for i, val := range vals {
		res[i] = []byte(val)
	}
	return res
}
