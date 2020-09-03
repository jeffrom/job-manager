// Package jobclient is an http client for interacting with job-manager server
// applications.
package jobclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"google.golang.org/protobuf/proto"

	apiv1 "github.com/jeffrom/job-manager/pkg/api/v1"
	"github.com/jeffrom/job-manager/pkg/job"
	"github.com/jeffrom/job-manager/pkg/schema"
)

type Interface interface {
	// Resource(name string) resource.Interface
	Ping(ctx context.Context) error

	// EnqueueJobs(ctx context.Context, jobs *job.Jobs) ([]string, error)
	// EnqueueJobsOpts(ctx context.Context, jobs *job.Jobs, opts EnqueueOpts) ([]string, error)
	EnqueueJob(ctx context.Context, job string, args ...interface{}) (string, error)
	// EnqueueJobOpts(ctx context.Context, jobData *job.Job, opts EnqueueOpts) error
	DequeueJobs(ctx context.Context, num int, job string, selectors ...string) (*job.Jobs, error)
	AckJob(ctx context.Context, id string, status job.Status) error
	AckJobOpts(ctx context.Context, id string, status job.Status, opts AckJobOpts) error
	// AckJobs(ctx context.Context, results *job.Results) error

	SaveQueue(ctx context.Context, name string, opts SaveQueueOpts) (*job.Queue, error)
	// SaveQueues(ctx context.Context, queue *job.Queues) error
	GetJob(ctx context.Context, id string) (*job.Job, error)
}

type providerFunc func(c *Client) *Client

type Client struct {
	addr   string
	cfg    *Config
	client *http.Client
}

func New(addr string, providers ...providerFunc) *Client {
	c := &Client{
		addr:   addr,
		cfg:    &ConfigDefaults,
		client: defaultClient(),
	}

	for _, provider := range providers {
		c = provider(c)
	}
	return c
}

func WithHTTPClient(client *http.Client) providerFunc {
	return func(c *Client) *Client {
		c.client = client
		return c
	}
}

func WithConfig(cfg *Config) providerFunc {
	return func(c *Client) *Client {
		c.cfg = cfg
		return c
	}
}

// func (c *Client) Resource(name string) resource.Interface {
// 	return nil
// }

func (c *Client) Ping(ctx context.Context) error {
	req, err := c.newRequest("GET", "/internal/ready", nil)
	if err != nil {
		return err
	}

	res, err := c.client.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("jobclient: ping failed with %d status", res.StatusCode)
	}
	return nil
}

func (c *Client) newRequest(method, uri string, body io.Reader) (*http.Request, error) {
	u := fmt.Sprintf("http://%s%s", c.addr, uri)
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("content-type", "application/protobuf")
	return req, nil
}

func (c *Client) newRequestProto(method, uri string, msg proto.Message) (*http.Request, error) {
	var r io.Reader
	if msg != nil {
		b, err := proto.Marshal(msg)
		if err != nil {
			return nil, err
		}
		r = bytes.NewReader(b)
	}
	return c.newRequest(method, uri, r)
}

func (c *Client) doRequest(ctx context.Context, req *http.Request, msg proto.Message) error {
	res, err := c.client.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}

	switch res.StatusCode {
	case http.StatusOK:
		if err := unmarshalProto(res, msg); err != nil {
			return err
		}
	case http.StatusNotFound:
		msg := &apiv1.GenericError{}
		if err := unmarshalProto(res, msg); err != nil {
			return err
		}
		return newNotFoundErrorProto(msg)
	case http.StatusBadRequest:
		msg := &apiv1.ValidationErrorResponse{}
		if err := unmarshalProto(res, msg); err != nil {
			return err
		}
		return schema.NewValidationErrorProto(msg)
	case http.StatusInternalServerError:
		return ErrInternal
	default:
		panic(fmt.Sprintf("unhandled exit code: %d %s", res.StatusCode, res.Status))
	}
	return nil
}

func unmarshalProto(res *http.Response, msg proto.Message) error {
	if msg == nil {
		return nil
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return proto.Unmarshal(b, msg)
}
