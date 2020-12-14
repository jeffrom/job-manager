package client

import (
	"context"
	"fmt"

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
