// Package bepostgres implements backend.Interface using Postgres.
package bepostgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	// "github.com/jackc/pgx/v4/pgxpool"
	// _ "github.com/jackc/pgx/v4/stdlib"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zerologadapter"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"

	"github.com/jeffrom/job-manager/mjob/label"
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

func (pg *Postgres) Close() error {
	if pg.db != nil {
		return pg.db.Close()
	}
	return nil
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
	connStr, err := registerConnConfig(dsn, pg.cfg.Logger.Logger, pg.cfg.Debug)
	if err != nil {
		return err
	}

	log.Debug().Str("dsn", dsn).Str("conn_str", connStr).Msg("connecting to postgres")
	db, err := sqlx.ConnectContext(ctx, "pgx", connStr)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(32)
	db.SetMaxIdleConns(32)
	db.SetConnMaxLifetime(15 * time.Minute)
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
		owner := ""
		table := ""
		if err := rows.Scan(&owner, &table); err != nil {
			return err
		}
		if table == "schema_migrations" {
			continue
		}
		tables = append(tables, table)
	}
	// fmt.Println(tables)

	for _, table := range tables {
		if _, err := pg.db.ExecContext(ctx, "TRUNCATE TABLE "+table+" CASCADE"); err != nil {
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
	defer stmt.Close()

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
	defer stmt.Close()

	for _, key := range stringsToBytea(keys) {
		if _, err := stmt.ExecContext(ctx, key); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}
	return nil
}

func registerConnConfig(dsn string, logger zerolog.Logger, debug bool) (string, error) {
	pgCfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		return "", err
	}
	pgCfg.Logger = zerologadapter.NewLogger(logger)
	if debug {
		pgCfg.LogLevel = pgx.LogLevelTrace
	} else {
		pgCfg.LogLevel = pgx.LogLevelError
	}

	return stdlib.RegisterConnConfig(pgCfg), nil
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

func sqlSelectors(sel *label.Selectors, joins, wheres []string, args []interface{}) ([]string, []string, []interface{}) {
	if sel.Len() == 0 {
		return joins, wheres, args
	}

	joins = append(joins, "JOIN queue_labels ON queues.name = queue_labels.queue")

	if names := sel.Names; len(names) > 0 {
		wheres = append(wheres, "queue_labels.name IN (?)")
		args = append(args, names)
	}
	if notnames := sel.NotNames; len(notnames) > 0 {
		wheres = append(wheres, "queue_labels.name NOT IN (?)")
		args = append(args, notnames)
	}
	if in := sel.In; len(in) > 0 {
		for k, v := range in {
			wheres = append(wheres, "queue_labels.name = ? AND queue_labels.value IN (?)")
			args = append(args, k, v)
		}
	}
	if notin := sel.NotIn; len(notin) > 0 {
		for k, v := range notin {
			wheres = append(wheres, "queue_labels.name = ? AND queue_labels.value NOT IN (?)")
			args = append(args, k, v)
		}
	}

	return joins, wheres, args
}
