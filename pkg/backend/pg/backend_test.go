package pg

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jeffrom/job-manager/pkg/backend"
	"github.com/jeffrom/job-manager/pkg/logger"
	"github.com/jeffrom/job-manager/pkg/testenv"
)

func TestBackendPostgres(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping because --short")
	}
	defaultCfg := backend.DefaultConfig
	defaultCfg.TestMode = true
	defaultCfg.Logger = logger.New(os.Stdout, false, true)
	cfg := Config{
		Config:   defaultCfg,
		Database: "job_manager_test",
	}
	ctx := context.Background()

	connStr, err := registerConnConfig("database=postgres", defaultCfg.Logger.Logger, true)
	if err != nil {
		t.Fatal(err)
	}

	conn, err := sql.Open("pgx", connStr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	findProjectRoot()
	resetDB(ctx, t, conn, cfg.Database)

	// make another conn for the migrations
	mconnStr, err := registerConnConfig("database=job_manager_test", defaultCfg.Logger.Logger, true)
	if err != nil {
		t.Fatal(err)
	}

	mconn, err := sql.Open("pgx", mconnStr)
	if err != nil {
		t.Fatal(err)
	}
	defer mconn.Close()
	runMigrations(ctx, t, mconn, cfg.Database)

	be := New(WithConfig(cfg))
	defer be.Close()
	testenv.BackendTest(testenv.BackendTestConfig{
		Type:    "postgres",
		Backend: be,
		Fail:    os.Getenv("CI") != "",
	})(t)
}

func resetDB(ctx context.Context, t testing.TB, conn *sql.DB, database string) {
	_, err := conn.ExecContext(ctx, "CREATE DATABASE "+database)
	if err != nil {
		t.Logf("ignoring error: %v", err)
	}
}

func runMigrations(ctx context.Context, t testing.TB, conn *sql.DB, database string) {
	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		t.Fatal(err)
	}

	root, err := findProjectRoot()
	if err != nil {
		t.Fatal(err)
	}

	p := filepath.Join(root, "pkg/backend/pg/migrations")
	m, err := migrate.NewWithDatabaseInstance("file://"+p, "postgres", driver)
	if err != nil {
		t.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatal(err)
	}
}

func findProjectRoot() (string, error) {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return "", errors.New("testenv: failed to get path of caller's file")
	}
	dir, _ := filepath.Split(file)
	for d := dir; d != "/"; d, _ = filepath.Split(filepath.Clean(d)) {
		gomodPath := filepath.Join(d, "go.mod")
		if _, err := os.Stat(gomodPath); err != nil {
			continue
		}
		return d, nil
	}
	return "", errors.New("didn't find project root")
}
