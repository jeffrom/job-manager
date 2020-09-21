package bepostgres

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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
	q := "INSERT INTO jobs (" + fields + ") VALUES (" + vals + ") RETURNING id"
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
	if opts == nil {
		opts = &resource.JobListParams{}
	}
	opts.NoUnclaimed = true
	opts.Statuses = []*resource.Status{resource.NewStatus(resource.StatusQueued), resource.NewStatus(resource.StatusFailed)}

	jobs, err := pg.ListJobs(ctx, limit, opts)
	if err != nil {
		return nil, err
	}

	now := internal.GetTimeProvider(ctx).Now().UTC()

	c, err := pg.getConn(ctx)
	if err != nil {
		return nil, err
	}

	q := "UPDATE jobs SET status = 'running', attempt = attempt+1, v = v+1, started_at = $1 WHERE id = $2 RETURNING *"
	stmt, err := sqlx.PreparexContext(ctx, c, q)
	if err != nil {
		return nil, err
	}

	resJobs := make([]*resource.Job, len(jobs.Jobs))
	for i, jb := range jobs.Jobs {
		fmt.Println("job:", jb)
		resJob := jb.Copy()
		fmt.Println("copy:", resJob)
		if err := stmt.GetContext(ctx, resJob, now, jb.ID); err != nil {
			return nil, err
		}
		fmt.Println("row:", resJob)

		resJobs[i] = resJob
	}
	jobs.Jobs = resJobs
	return jobs, nil
}

func (pg *Postgres) AckJobs(ctx context.Context, results *resource.Acks) error {
	return nil
}

func (pg *Postgres) GetJobByID(ctx context.Context, id string) (*resource.Job, error) {
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
	now := internal.GetTimeProvider(ctx).Now().UTC()

	// XXX for claims (to get the claim duration) and selectors (for
	// queue_labels), we need a queue map. for claims the version matters, for
	// labels it doesn't.

	q := fmt.Sprintf("%s FROM jobs LEFT JOIN queues ON jobs.queue_id = queues.id", jobFields)
	usingLabels := false
	var joins []string
	var wheres []string
	var args []interface{}

	if len(opts.Names) > 0 {
		wheres = append(wheres, "queues.name IN (?)")
		args = append(args, opts.Names)
	}
	if len(opts.Statuses) > 0 {
		wheres = append(wheres, "jobs.status IN (?)")
		args = append(args, resource.StatusStrings(opts.Statuses...))
	}
	if opts.Selectors.Len() > 0 {
		joins, wheres, args = sqlSelectors(opts.Selectors, joins, wheres, args)
		usingLabels = true
	}
	if opts.NoUnclaimed || len(opts.Claims) > 0 {
		joins = append(joins, "LEFT JOIN job_claims ON jobs.id = job_claims.job_id")
		joins = append(joins, "LEFT JOIN (SELECT DISTINCT ON (job_id) job_id, completed_at AS completed_at FROM job_results ORDER BY job_id, completed_at DESC) AS last_attempt ON jobs.id = last_attempt.job_id")
	}
	if opts.NoUnclaimed && len(opts.Claims) == 0 {
		wheres = append(wheres, "GREATEST(jobs.enqueued_at, last_attempt.completed_at) + (queues.claim_duration * INTERVAL '1 microsecond') <= ?")
		args = append(args, now)
	}
	if len(opts.Claims) > 0 {
		for name, vals := range opts.Claims {
			wheres = append(wheres, "(job_claims.name = ? AND job_claims.value IN (?)) OR (GREATEST(jobs.enqueued_at, last_attempt.completed_at) + (queues.claim_duration * INTERVAL '1 microsecond') <= ?)")
			args = append(args, name, vals, now)
		}
	}

	if usingLabels {
		q = "SELECT DISTINCT ON (id) " + q
	} else {
		q = "SELECT " + q
	}

	if len(joins) > 0 {
		q += " " + strings.Join(joins, " ")
	}
	if len(wheres) > 0 {
		q += " WHERE " + strings.Join(wheres, " AND ")
	}

	q, args, err = sqlx.In(q, args...)
	if err != nil {
		return nil, err
	}

	var rows []*resource.Job
	if err := sqlx.SelectContext(ctx, c, &rows, c.Rebind(q), args...); err != nil {
		return nil, err
	}
	for _, row := range rows {
		if err := row.Populate(); err != nil {
			return nil, err
		}
	}

	// XXX if we're using labels have to query queue labels here to handle the
	// !label selector :'(

	if err := annotateJobs(ctx, c, rows); err != nil {
		return nil, err
	}

	return &resource.Jobs{Jobs: rows}, nil
}

func annotateJobs(ctx context.Context, c sqlxer, jobs []*resource.Job) error {
	if len(jobs) == 0 {
		return nil
	}
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
