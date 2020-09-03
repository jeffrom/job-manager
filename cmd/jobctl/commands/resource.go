package commands

import (
	"github.com/jeffrom/job-manager/jobclient"
)

type resourceOpts struct{}

type cliResource struct {
	cfg    jobclient.Config
	opts   *resourceOpts
	client jobclient.Interface
}

func newResource(cfg jobclient.Config, opts *resourceOpts) *cliResource {
	if opts == nil {
		opts = &resourceOpts{}
	}

	return &cliResource{
		cfg:    cfg,
		opts:   opts,
		client: jobclient.New(cfg.Addr, jobclient.WithConfig(cfg)),
	}
}
