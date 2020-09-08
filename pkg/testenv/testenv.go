// Package testenv contains helpers for setting up job-manager test
// environments.
package testenv

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/jeffrom/job-manager/jobclient"
	"github.com/jeffrom/job-manager/pkg/backend"
	"github.com/jeffrom/job-manager/pkg/web"
	"github.com/jeffrom/job-manager/pkg/web/middleware"
)

func NewTestControllerServer(t testing.TB, cfg middleware.Config, be backend.Interface) *httptest.Server {
	t.Helper()
	cfg.ResetLogOutput(testLogOutput(t))
	if cfg.Backend == "" {
		cfg.Backend = "memory"
	}
	h, err := web.NewControllerRouter(cfg, be)
	die(t, err)

	srv := httptest.NewUnstartedServer(h)
	t.Logf("Started job-controller server with backend %T at address: %s", be, srv.Listener.Addr())
	return srv
}

func NewTestClient(t testing.TB, srv *httptest.Server) *jobclient.Client {
	return jobclient.New(srv.Listener.Addr().String())
}

func die(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

type testLogger struct {
	t testing.TB
}

func (t *testLogger) Write(b []byte) (int, error) {
	t.t.Log(string(b))
	return len(b), nil
}

func testLogOutput(t testing.TB) io.Writer {
	return &testLogger{t: t}
}
