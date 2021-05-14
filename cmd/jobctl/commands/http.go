package commands

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/jeffrom/job-manager/release"
)

type roundTripper struct {
	*http.Transport
}

func newRoundTripper(transport *http.Transport) *roundTripper {
	return &roundTripper{Transport: transport}
}

func (rt *roundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("User-agent", fmt.Sprintf("jobctl/%s", release.Version))
	return rt.Transport.RoundTrip(r)
}

var httpClient = &http.Client{
	Timeout: 15 * time.Second,
	Transport: newRoundTripper(&http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 5 * time.Second,
	}),
}
