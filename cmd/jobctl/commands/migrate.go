package commands

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/mjob/client"
)

type migrateCmd struct {
	*cobra.Command
}

func newMigrateCmd(cfg *client.Config) *migrateCmd {
	c := &migrateCmd{
		Command: &cobra.Command{
			Use:  "migrate",
			Args: cobra.NoArgs,
		},
	}

	return c
}

func (c *migrateCmd) Cmd() *cobra.Command { return c.Command }
func (c *migrateCmd) Execute(ctx context.Context, cfg *client.Config, cmd *cobra.Command, args []string) error {
	// TODO flags for timeouts, etc
	client := &http.Client{
		Timeout: 5 * time.Minute,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}
	uri := fmt.Sprintf("http://%s/api/v1/backend/migrate", cfg.Host)
	req, err := http.NewRequest("POST", uri, nil)
	if err != nil {
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode >= 300 {
		defer res.Body.Close()
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("got %d error: %s", res.StatusCode, b)
	}
	return nil
}
