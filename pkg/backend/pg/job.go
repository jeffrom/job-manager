package pg

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/jeffrom/job-manager/mjob/resource"
	"github.com/jeffrom/job-manager/pkg/backend"
	"github.com/jeffrom/job-manager/pkg/internal"
)

const (
	jobFields = "jobs.id, jobs.v, queues.name, queues.v AS queue_v, queues.duration, queues.backoff_initial_duration, queues.backoff_max_duration, queues.backoff_factor, attempt, status, args, data, enqueued_at, started_at"
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
	for _, q := range res.Queues {
		if q.Blocked {
			return nil, backend.ErrBlocked
		}
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
		if jb.Data != nil {
			jb.DataRaw = jb.Data.DataRaw
		}
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

	claimQ := "INSERT INTO job_claims (job_id, name, value) VALUES (?, ?, ?)"
	insertClaims, err := c.PrepareContext(ctx, c.Rebind(claimQ))
	if err != nil {
		return nil, err
	}
	defer insertClaims.Close()

	for _, jb := range jobs.Jobs {
		row := stmt.QueryRowContext(ctx, jb)
		var id int64
		if err := row.Scan(&id); err != nil {
			return nil, err
		}
		jb.ID = strconv.FormatInt(id, 10)

		if jb.Data != nil && len(jb.Data.Claims) > 0 {
			for k, vals := range jb.Data.Claims {
				for _, v := range vals {
					if _, err := insertClaims.ExecContext(ctx, jb.ID, k, v); err != nil {
						return nil, err
					}
				}
			}
		}
	}

	return jobs, nil
}

func (pg *Postgres) DequeueJobs(ctx context.Context, limit int, opts *resource.JobListParams) (*resource.Jobs, error) {
	if opts == nil {
		opts = &resource.JobListParams{}
	}
	opts.NoUnclaimed = true
	opts.NoPaused = true
	opts.Statuses = []*resource.Status{resource.NewStatus(resource.StatusQueued), resource.NewStatus(resource.StatusFailed)}

	// NOTE annotate jobs after the insert loop below if we want to include
	// job_results when start-while-running happens.
	jobs, err := pg.listJobs(ctx, limit, opts, true)
	if err != nil {
		return nil, err
	}

	now := internal.GetTimeProvider(ctx).Now().UTC()

	c, err := pg.getConn(ctx)
	if err != nil {
		return nil, err
	}

	// NOTE we could insert a job_results row here, but because we don't
	// actually know what happenned, it makes more sense not to try to pretend
	// that it did. Plus we don't know when the last job failed, so we could
	// mess up backoff.

	updateJobQ := "UPDATE jobs SET status = 'running', attempt = attempt+1, v = v+1, started_at = $1 WHERE id = $2 RETURNING *"
	stmt, err := sqlx.PreparexContext(ctx, c, updateJobQ)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	resJobs := make([]*resource.Job, len(jobs.Jobs))
	for i, jb := range jobs.Jobs {
		resJob := jb.Copy()
		if err := stmt.GetContext(ctx, resJob, now, jb.ID); err != nil {
			return nil, err
		}

		resJobs[i] = resJob
	}
	jobs.Jobs = resJobs
	return jobs, nil
}

func (pg *Postgres) AckJobs(ctx context.Context, results *resource.Acks) error {
	now := internal.GetTimeProvider(ctx).Now().UTC()
	c, err := pg.getConn(ctx)
	if err != nil {
		return err
	}

	updateQ := "UPDATE jobs SET status = ?, v = v+1, completed_at = ? WHERE id = ? RETURNING *"
	update, err := sqlx.PreparexContext(ctx, c, c.Rebind(updateQ))
	if err != nil {
		return err
	}
	defer update.Close()

	insert, err := prepareInsertJobResults(ctx, c)
	if err != nil {
		return err
	}
	defer insert.Close()

	for _, ack := range results.Acks {
		jb := &resource.Job{}
		status := ack.Status.String()
		if err := update.GetContext(ctx, jb, status, now, ack.JobID); err != nil {
			return err
		}

		// for cancelled or invalid, we dont presume the job has already started.
		// TODO we should probably be smarter about seeing if the job has
		// started or not and fix the update query accordingly, or may be
		// smarter to allow job_results.started_at to be null, which would
		// indicate that the job was never started, but we still wanted to
		// attach some data to it.
		if *ack.Status == resource.StatusCancelled || *ack.Status == resource.StatusInvalid {
			continue
		}

		var data []byte
		if ack.Data != nil {
			data, err = json.Marshal(ack.Data)
			if err != nil {
				return err
			}
		}
		if _, err := insert.ExecContext(ctx, ack.JobID, ack.Status, data, ack.Error, jb.StartedAt, now); err != nil {
			return err
		}
	}
	return nil
}

func (pg *Postgres) GetJobByID(ctx context.Context, id string, opts *resource.GetByIDOpts) (*resource.Job, error) {
	if opts == nil {
		opts = &resource.GetByIDOpts{}
	}
	c, err := pg.getConn(ctx)
	if err != nil {
		return nil, err
	}

	jb := &resource.Job{}
	q := fmt.Sprintf("SELECT %s FROM jobs LEFT JOIN queues ON jobs.queue_id = queues.id WHERE jobs.id = $1", jobFields)
	if err := sqlx.GetContext(ctx, c, jb, q, id); err != nil {
		return nil, err
	}

	if err := annotateJobs(ctx, c, opts.Includes, []*resource.Job{jb}); err != nil {
		return nil, err
	}
	// fmt.Printf("GetJobByID: Args: %q\n", jb.ArgsRaw)
	return jb, nil
}

func (pg *Postgres) ListJobs(ctx context.Context, limit int, opts *resource.JobListParams) (*resource.Jobs, error) {
	return pg.listJobs(ctx, limit, opts, false)
}

func (pg *Postgres) listJobs(ctx context.Context, limit int, opts *resource.JobListParams, forDequeue bool) (*resource.Jobs, error) {
	if opts == nil {
		opts = &resource.JobListParams{}
	}
	// fmt.Printf("ASDF: %+v\n", opts)
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

	if len(opts.Queues) > 0 {
		wheres = append(wheres, "queues.name IN (?)")
		args = append(args, opts.Queues)
	}
	if len(opts.Statuses) > 0 {
		if forDequeue {
			// handle case where job duration has elapsed and the reaper hasn't
			// updated status yet
			wheres = append(wheres, "(jobs.status IN (?) OR (jobs.status = 'running' AND ((jobs.started_at IS NOT NULL AND ? > jobs.started_at + ((queues.duration / 1000) * INTERVAL '1 microsecond')) OR (jobs.started_at IS NOT NULL AND ? > jobs.started_at + ((queues.duration / 1000) * INTERVAL '1 microsecond')))))")
			args = append(args, resource.StatusStrings(opts.Statuses...), now, now)
		} else {
			wheres = append(wheres, "jobs.status IN (?)")
			args = append(args, resource.StatusStrings(opts.Statuses...))
		}
	}
	if opts.Selectors.Len() > 0 {
		joins, wheres, args = sqlSelectors(opts.Selectors, joins, wheres, args)
		usingLabels = true
	}
	if opts.NoUnclaimed || len(opts.Claims) > 0 {
		joins = append(joins, "LEFT JOIN job_claims ON jobs.id = job_claims.job_id")
		// joins = append(joins, "LEFT JOIN (SELECT DISTINCT ON (job_id) job_id, started_at AS last_attempt_started_at, completed_at AS completed_at FROM job_results ORDER BY job_id, id DESC) AS last_attempt ON jobs.id = last_attempt.job_id")
	}
	if opts.NoUnclaimed && len(opts.Claims) == 0 {
		wheres = append(wheres, "(job_claims.job_id IS NULL OR (GREATEST(jobs.enqueued_at, jobs.completed_at) + ((queues.claim_duration / 1000) * INTERVAL '1 microsecond') <= ?))")
		args = append(args, now)
	}
	if len(opts.Claims) > 0 {
		for name, vals := range opts.Claims {
			wheres = append(wheres, "((job_claims.name = ? AND job_claims.value IN (?)) OR (GREATEST(jobs.enqueued_at, jobs.completed_at) + ((queues.claim_duration / 1000) * INTERVAL '1 microsecond') <= ?))")
			args = append(args, name, vals, now)
		}
	}
	if !opts.EnqueuedSince.IsZero() {
		wheres = append(wheres, "(jobs.enqueued_at >= ?)")
		args = append(args, opts.EnqueuedSince.UTC())
	}
	if !opts.EnqueuedUntil.IsZero() {
		wheres = append(wheres, "(jobs.enqueued_at < ?)")
		args = append(args, opts.EnqueuedSince.UTC())
	}
	if opts.Page != nil && opts.Page.LastID != "" {
		op := "<"
		if forDequeue {
			op = ">"
		}
		wheres = append(wheres, fmt.Sprintf("(jobs.id %s ?)", op))
		args = append(args, opts.Page.LastID)
	}
	if opts.NoPaused {
		wheres = append(wheres, "(queues.paused != true OR queues.unpaused = true)")
	}
	if forDequeue {
		wheres = append(wheres, "(jobs.attempt <= queues.retries)")
		wheres = append(wheres, "(queues.backoff_initial_duration = 0 OR queues.backoff_factor = 0 OR jobs.completed_at IS NULL OR (? > jobs.completed_at + (LEAST(queues.backoff_max_duration, (queues.backoff_initial_duration * (jobs.attempt ^ queues.backoff_factor)) / 1000) * INTERVAL '1 microsecond')))")
		args = append(args, now)
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

	if forDequeue {
		q += " ORDER BY id ASC"
	} else {
		q += " ORDER BY id DESC"
	}

	q += fmt.Sprintf(" LIMIT %d", limit)
	if forDequeue {
		q += " FOR UPDATE OF jobs SKIP LOCKED"
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
		if len(row.DataRaw) > 0 {
			row.Data = &resource.JobData{DataRaw: row.DataRaw}
		}
	}

	// XXX if we're using labels have to query queue labels here to handle the
	// !label selector :'(

	if err := annotateJobs(ctx, c, opts.Includes, rows); err != nil {
		return nil, err
	}

	return &resource.Jobs{Jobs: rows}, nil

}

func annotateJobs(ctx context.Context, c sqlxer, includes []string, jobs []*resource.Job) error {
	if len(jobs) == 0 {
		return nil
	}
	if len(includes) == 0 {
		return nil
	}

	ids := make([]string, len(jobs))
	jobmap := make(map[string]*resource.Job)
	for i, jb := range jobs {
		ids[i] = jb.ID
		jobmap[jb.ID] = jb
	}
	incMap := makeIncludeMap(includes)

	if incMap["checkin"] {
		q, args, err := sqlx.In("SELECT * FROM job_checkins WHERE job_id in (?)", ids)
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
	}

	if incMap["result"] {
		q, args, err := sqlx.In("SELECT * FROM job_results WHERE job_id in (?)", ids)
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
	}
	return nil
}

func prepareInsertJobResults(ctx context.Context, c sqlxer) (*sqlx.Stmt, error) {
	sql := "INSERT INTO job_results (job_id, status, data, error, started_at, completed_at) VALUES (?, ?, ?, ?, ?, ?)"
	q, err := sqlx.PreparexContext(ctx, c, c.Rebind(sql))
	return q, err
}
