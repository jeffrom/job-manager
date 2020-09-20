package bepostgres

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/jeffrom/job-manager/pkg/internal"
	"github.com/jeffrom/job-manager/pkg/resource"
	"github.com/jmoiron/sqlx"
)

const (
	jobFields = "jobs.id, jobs.v, queues.name, queues.v AS queue_v, attempt, status, args, data, enqueued_at, started_at"
)

func (pg *Postgres) EnqueueJobs(ctx context.Context, jobs *resource.Jobs) (*resource.Jobs, error) {
	c, err := pg.getConn(ctx)
	if err != nil {
		return nil, err
	}
	res, err := pg.ListQueues(ctx, &resource.QueueListParams{Names: jobs.Queues()})
	if err != nil {
		return nil, err
	}
	queues := res.ToMap()

	// var uniquenessKeys []string
	now := internal.GetTimeProvider(ctx).Now().UTC()
	for _, jb := range jobs.Jobs {
		q := queues[jb.Name]
		jb.EnqueuedAt = now
		jb.QueueID = q.ID
		jb.Version = resource.NewVersion(1)
		jb.QueueVersion = q.Version
		jb.Status = resource.NewStatus(resource.StatusQueued)
	}

	fields, vals := sqlFields(
		"v", "queue_id",
		"attempt", "status",
		"args", "data", "enqueued_at",
	)
	q := "INSERT INTO jobs (" + fields + ", root_id) VALUES (" + vals + ", 0) RETURNING id"
	stmt, err := c.PrepareNamedContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	for _, jb := range jobs.Jobs {
		row := stmt.QueryRowContext(ctx, jb)
		var id int64
		if err := row.Scan(&id); err != nil {
			return nil, err
		}
		jb.ID = strconv.FormatInt(id, 10)
	}

	return jobs, nil
}

func (pg *Postgres) DequeueJobs(ctx context.Context, limit int, opts *resource.JobListParams) (*resource.Jobs, error) {
	return nil, nil
}

func (pg *Postgres) AckJobs(ctx context.Context, results *resource.Acks) error {
	return nil
}

func (pg *Postgres) GetJobByID(ctx context.Context, id string) (*resource.Job, error) {
	// id, err := strconv.ParseInt(idstr, 10, 64)
	// if err != nil {
	// 	return nil, err
	// }
	c, err := pg.getConn(ctx)
	if err != nil {
		return nil, err
	}

	jb := &resource.Job{}
	q := fmt.Sprintf("SELECT %s FROM jobs LEFT JOIN queues ON jobs.queue_id = queues.id WHERE jobs.id = $1", jobFields)
	if err := sqlx.GetContext(ctx, c, jb, q, id); err != nil {
		return nil, err
	}

	if err := annotateJobs(ctx, c, []*resource.Job{jb}); err != nil {
		return nil, err
	}
	return jb, nil
}

func (pg *Postgres) ListJobs(ctx context.Context, limit int, opts *resource.JobListParams) (*resource.Jobs, error) {
	if opts == nil {
		opts = &resource.JobListParams{}
	}
	c, err := pg.getConn(ctx)
	if err != nil {
		return nil, err
	}

	q := fmt.Sprintf("SELECT %s FROM jobs LEFT JOIN queues ON jobs.queue_id = queues.id", jobFields)
	// var froms []string
	var wheres []string
	var args []interface{}

	if len(opts.Names) > 0 {
		wheres = append(wheres, "queues.name IN (?)")
		args = append(args, opts.Names)
	}
	if len(opts.Claims) > 0 {

	}
	if len(opts.Statuses) > 0 {

	}
	if len(opts.Claims) > 0 {

	}

	var rows []*resource.Job
	if err := sqlx.SelectContext(ctx, c, &rows, q); err != nil {
		return nil, err
	}
	return &resource.Jobs{Jobs: rows}, nil
}

func uniquenessKeyFromArgs(args []interface{}) (string, error) {
	b, err := json.Marshal(args)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return string(sum[:]), nil
}

func annotateJobs(ctx context.Context, c sqlxer, jobs []*resource.Job) error {
	ids := make([]string, len(jobs))
	jobmap := make(map[string]*resource.Job)
	for i, jb := range jobs {
		ids[i] = jb.ID
		jobmap[jb.ID] = jb
	}

	q, args, err := sqlx.In("SELECT * FROM job_checkins WHERE id in (?)", ids)
	if err != nil {
		return err
	}

	var checkins []*resource.JobCheckin
	if err := sqlx.SelectContext(ctx, c, &checkins, c.Rebind(q), args...); err != nil {
		return err
	}
	for _, row := range checkins {
		jb := jobmap[row.JobID]
		jb.Checkins = append(jb.Checkins, row)
	}

	q, args, err = sqlx.In("SELECT * FROM job_results WHERE id in (?)", ids)
	if err != nil {
		return err
	}

	var results []*resource.JobResult
	if err := sqlx.SelectContext(ctx, c, &results, c.Rebind(q), args...); err != nil {
		return err
	}
	for _, row := range results {
		jb := jobmap[row.JobID]
		jb.Results = append(jb.Results, row)
	}
	return nil
}
