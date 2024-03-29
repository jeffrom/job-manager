// Package client contains the base job-manager http client.
package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/protobuf/proto"

	apiv1 "github.com/jeffrom/job-manager/mjob/api/v1"
	"github.com/jeffrom/job-manager/mjob/internal"
	"github.com/jeffrom/job-manager/mjob/querystring"
	"github.com/jeffrom/job-manager/mjob/resource"
)

type Interface interface {
	// Resource(name string) resource.Interface
	Ping(ctx context.Context) error

	// consumer rpcs
	// EnqueueJobs(ctx context.Context, jobs *resource.Jobs) ([]string, error)
	// EnqueueJobsOpts(ctx context.Context, jobs *resource.Jobs, opts EnqueueOpts) ([]string, error)
	EnqueueJob(ctx context.Context, job string, args ...interface{}) (string, error)
	EnqueueJobOpts(ctx context.Context, job string, opts EnqueueOpts, args ...interface{}) (string, error)
	DequeueJobs(ctx context.Context, num int, id string) (*resource.Jobs, error)
	DequeueJobsOpts(ctx context.Context, num int, opts DequeueOpts) (*resource.Jobs, error)
	AckJob(ctx context.Context, id string, status resource.Status) error
	AckJobOpts(ctx context.Context, id string, status resource.Status, opts AckJobOpts) error
	// AckJobs(ctx context.Context, results *resource.Results) error

	ListJobs(ctx context.Context, opts ListJobsOpts) (*resource.Jobs, error)
	GetJob(ctx context.Context, id string, opts *GetJobOpts) (*resource.Job, error)

	SaveQueue(ctx context.Context, name string, opts SaveQueueOpts) (*resource.Queue, error)
	// SaveQueues(ctx context.Context, queue *resource.Queues) error
	ListQueues(ctx context.Context, opts ListQueuesOpts) (*resource.Queues, error)
	GetQueue(ctx context.Context, id string) (*resource.Queue, error)
	DeleteQueue(ctx context.Context, queue string) error
	PauseQueue(ctx context.Context, queue string) error
	UnpauseQueue(ctx context.Context, queue string) error
	BlockQueue(ctx context.Context, queue string) error
	UnblockQueue(ctx context.Context, queue string) error
}

type Provider func(c *Client) *Client

type Client struct {
	addr   string
	cfg    *Config
	client *http.Client
}

func New(addr string, providers ...Provider) *Client {
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

func WithHTTPClient(client *http.Client) Provider {
	return func(c *Client) *Client {
		c.client = client
		return c
	}
}

func WithConfig(cfg *Config) Provider {
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
		return fmt.Errorf("client: ping failed with %d status", res.StatusCode)
	}
	return nil
}

func (c *Client) Stats(ctx context.Context, queue string) (*resource.Stats, error) {
	uri := "/api/v1/stats"
	if queue != "" {
		uri += "/" + queue
	}
	req, err := c.newRequestProto(ctx, "GET", uri, nil)
	if err != nil {
		return nil, err
	}
	res := &apiv1.StatsResponse{}
	if err := c.doRequest(ctx, req, res); err != nil {
		return nil, err
	}
	stats := &resource.Stats{
		Queued:           res.Queued,
		Running:          res.Running,
		Complete:         res.Complete,
		Dead:             res.Dead,
		Cancelled:        res.Cancelled,
		Invalid:          res.Invalid,
		Failed:           res.Failed,
		LongestUnstarted: res.LongestUnstartedSecs,
	}
	return stats, nil
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

func (c *Client) newRequestProto(ctx context.Context, method, uri string, msg proto.Message) (*http.Request, error) {
	var r io.Reader
	if msg != nil {
		if method == "GET" {
			vals, err := querystring.Values(msg)
			if err != nil {
				return nil, err
			}
			uri += "?" + vals.Encode()
		} else {
			b, err := proto.Marshal(msg)
			if err != nil {
				return nil, err
			}
			r = bytes.NewReader(b)
		}
	}
	// fmt.Printf("uri: %q\n", uri)
	req, err := c.newRequest(method, uri, r)
	// if method == "GET" {
	// 	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	// }
	// b, _ := httputil.DumpRequest(req, false)
	// fmt.Println(string(b))
	if mockNow := internal.GetMockTime(ctx); mockNow != nil {
		timeStr := fmt.Sprint(mockNow.Unix())
		req.Header.Set("fake-time", timeStr)
	}

	return req, err
}

func (c *Client) doRequest(ctx context.Context, req *http.Request, msg proto.Message) error {
	res, err := c.client.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}

	// b, _ := httputil.DumpResponse(res, false)
	// fmt.Println(string(b))

	switch code := res.StatusCode; {
	case code >= 200 && code < 300:
		if err := unmarshalProto(res, msg); err != nil {
			return err
		}
	default:
		msg := &apiv1.GenericError{}
		if err := unmarshalProto(res, msg); err != nil {
			return err
		}
		return newResourceErrorFromMessage(msg)
	}

	return nil
}

func unmarshalProto(res *http.Response, msg proto.Message) error {
	if msg == nil {
		return nil
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return proto.Unmarshal(b, msg)
}
