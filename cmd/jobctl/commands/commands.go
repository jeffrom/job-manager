// Package commands contains jobctl's cobra commands.
package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/jeffrom/job-manager/jobclient"
)

func Execute() {
	cfg := &jobclient.Config{}
	cmd := newRootCmd(cfg)
	ctx := context.Background()
	if err := cmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
