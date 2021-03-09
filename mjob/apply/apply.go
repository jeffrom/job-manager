// Package apply saves a yaml manifest to a job-manager server.
package apply

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/ghodss/yaml"

	"github.com/jeffrom/job-manager/mjob/client"
	"github.com/jeffrom/job-manager/mjob/resource"
	"github.com/jeffrom/job-manager/mjob/schema"
)

var docSep = []byte("\n---\n")

func Path(ctx context.Context, c client.Interface, p string) error {
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}

	docs := bytes.Split(b, docSep)
	for _, part := range docs {
		q := &resource.Queue{}
		if err := yaml.Unmarshal(part, q); err != nil {
			return err
		}
		if err := setSchema(ctx, part, q); err != nil {
			return err
		}
		res, err := applyRequest(ctx, c, q)
		if err != nil {
			return err
		}
		fmt.Printf("<- %s\n", res)
	}
	return nil
}

func applyRequest(ctx context.Context, c client.Interface, q *resource.Queue) (*resource.Queue, error) {
	prev, err := c.GetQueue(ctx, q.Name)
	if err != nil && !client.IsNotFound(err) {
		return nil, err
	}
	v := ""
	if prev != nil {
		v = prev.Version.Strict()
	}

	opts := toSaveOpts(q, v)
	return c.SaveQueue(ctx, q.Name, opts)
}

func toSaveOpts(q *resource.Queue, v string) client.SaveQueueOpts {
	return client.SaveQueueOpts{
		MaxRetries:      q.Retries,
		JobDuration:     time.Duration(q.Duration),
		CheckinDuration: time.Duration(q.CheckinDuration),
		ClaimDuration:   time.Duration(q.ClaimDuration),
		Labels:          q.Labels,
		Schema:          q.SchemaRaw,
		Unique:          q.Unique,
		Version:         v,
		BackoffInitial:  time.Duration(q.BackoffInitial),
		BackoffMax:      time.Duration(q.BackoffMax),
		BackoffFactor:   q.BackoffFactor,
	}
}

func setSchema(ctx context.Context, part []byte, q *resource.Queue) error {
	d := struct {
		Schema *schema.Schema `json:"schema"`
	}{}
	if err := yaml.Unmarshal(part, &d); err != nil {
		return err
	}
	if d.Schema == nil {
		return nil
	}
	b, err := json.Marshal(d.Schema)
	if err != nil {
		return err
	}
	if err := schema.ValidateSchema(ctx, b); err != nil {
		return err
	}
	q.SchemaRaw = b
	return nil
}
