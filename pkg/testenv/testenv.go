// Package testenv contains helpers for setting up job-manager test
// environments.
package testenv

import (
	"net/http/httptest"
	"testing"

	"github.com/jeffrom/job-manager/jobclient"
	"github.com/jeffrom/job-manager/pkg/backend"
	"github.com/jeffrom/job-manager/pkg/web"
)

func NewTestControllerServer(t testing.TB, cfg web.Config) *httptest.Server {
	t.Helper()
	be := cfg.GetBackend()
	if be == nil {
		be = backend.NewMemory()
	}
	h, err := web.NewControllerRouter(be)
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
