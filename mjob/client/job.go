package client

import (
	"context"
	"fmt"

	apiv1 "github.com/jeffrom/job-manager/mjob/api/v1"
	"github.com/jeffrom/job-manager/mjob/label"
	"github.com/jeffrom/job-manager/mjob/resource"
	jobv1 "github.com/jeffrom/job-manager/mjob/resource/job/v1"
)

func (c *Client) GetJob(ctx context.Context, id string) (*resource.Job, error) {
	uri := fmt.Sprintf("/api/v1/jobs/%s", id)
	req, err := c.newRequestProto(ctx, "GET", uri, nil)
	if err != nil {
		return nil, err
	}

	jobData := &jobv1.Job{}
	if err := c.doRequest(ctx, req, jobData); err != nil {
		return nil, err
	}
	var claims label.Claims
	if jobData.Data != nil && len(jobData.Data.Claims) > 0 {
		claims, err = label.ParseClaims(jobData.Data.Claims)
		if err != nil {
			return nil, err
		}
	}
	return jobv1.NewJobFromProto(jobData, claims), nil
}

type ListJobsOpts struct {
	Queues    []string
	Selectors []string
	Statuses  []resource.Status
	Page      *resource.Pagination
}

func (c *Client) ListJobs(ctx context.Context, opts ListJobsOpts) (*resource.Jobs, error) {
	page := apiv1.PaginationToProto(opts.Page)
	params := &apiv1.ListJobsRequest{
		Queue:    opts.Queues,
		Selector: opts.Selectors,
		Status:   statusStrings(opts.Statuses),
		Page:     page,
	}
	uri := "/api/v1/jobs"
	req, err := c.newRequestProto(ctx, "GET", uri, params)
	if err != nil {
		return nil, err
	}

	resp := &apiv1.ListJobsResponse{}
	if err := c.doRequest(ctx, req, resp); err != nil {
		return nil, err
	}
	jobs, err := jobv1.NewJobsFromProto(resp.Items)
	if err != nil {
		return nil, err
	}
	return &resource.Jobs{Jobs: jobs}, nil
}

func statusStrings(statuses []resource.Status) []string {
	res := make([]string, len(statuses))
	for i, st := range statuses {
		res[i] = st.String()
	}
	return res
}
