// Package mem is a minimal reference implementation of backend.Interface
// suitable for use in tests.
package mem

import (
	"context"
	"math"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/jeffrom/job-manager/mjob/label"
	"github.com/jeffrom/job-manager/mjob/resource"
	"github.com/jeffrom/job-manager/pkg/backend"
	"github.com/jeffrom/job-manager/pkg/internal"
)

// Memory is an in-memory backend intended to be a reference implementation
// used for testing. It is not safe to use in production.
type Memory struct {
	mu         sync.Mutex
	queues     map[string]*resource.Queue
	jobs       map[string]*resource.Job
	uniqueness map[string]bool
}

func New() *Memory {
	return &Memory{
		queues:     make(map[string]*resource.Queue),
		jobs:       make(map[string]*resource.Job),
		uniqueness: make(map[string]bool),
	}
}

func (m *Memory) Ping(ctx context.Context) error {
	return nil
}

func (m *Memory) Reset(ctx context.Context) error {
	m.queues = make(map[string]*resource.Queue)
	m.jobs = make(map[string]*resource.Job)
	m.uniqueness = make(map[string]bool)
	return nil
}

func (m *Memory) GetQueue(ctx context.Context, name string, opts *resource.GetByIDOpts) (*resource.Queue, error) {
	if cfg, ok := m.queues[name]; ok {
		return cfg, nil
	}
	return nil, backend.ErrNotFound
}

func (m *Memory) SaveQueue(ctx context.Context, queue *resource.Queue) (*resource.Queue, error) {
	if queue.Version == nil {
		queue.Version = resource.NewVersion(0)
	}
	// if it already exists and no version was supplied, or if version was
	// supplied but they don't match, return conflict
	prev, ok := m.queues[queue.Name]
	// fmt.Printf("prev: %+v, found: %v\n", prev, ok)
	if ok {
		if queue.Version.Raw() == 0 || queue.Version.Raw() != prev.Version.Raw() {
			return nil, &backend.VersionConflictError{
				Resource:   "queue",
				ResourceID: queue.Name,
				Prev:       prev.Version.String(),
				Curr:       queue.Version.String(),
			}
		}
	}
	// fmt.Printf("prev: %+v, curr: %+v\n", prev, queue)
	if prev == nil || !queue.EqualAttrs(prev) {
		now := internal.GetTimeProvider(ctx).Now().UTC()
		if queue.CreatedAt.IsZero() {
			queue.CreatedAt = now
		}
		queue.UpdatedAt = now
		queue.Version.Inc()
	}
	m.queues[queue.Name] = queue.Copy()
	return queue, nil
}

func (m *Memory) ListQueues(ctx context.Context, opts *resource.QueueListParams) (*resource.Queues, error) {
	if opts == nil {
		opts = &resource.QueueListParams{}
	}
	sels := opts.Selectors
	// sels, err := label.ParseSelectorStringArray(opts.Selectors)
	// if err != nil {
	// 	return nil, err
	// }

	queues := &resource.Queues{}
	for _, queue := range m.queues {
		if m.filterQueue(queue, opts.Names, sels) {
			continue
		}
		queues.Queues = append(queues.Queues, queue)
	}
	return queues, nil
}

func (m *Memory) filterQueue(queue *resource.Queue, names []string, sels *label.Selectors) bool {
	if len(names) > 0 && !valIn(queue.Name, names) {
		return true
	}
	return !sels.Match(queue.Labels)
}

func (m *Memory) EnqueueJobs(ctx context.Context, jobArgs *resource.Jobs) (*resource.Jobs, error) {
	now := internal.GetTimeProvider(ctx).Now().UTC()
	for _, jobArg := range jobArgs.Jobs {
		queue, err := m.GetQueue(ctx, jobArg.Name, nil)
		if err != nil {
			return nil, err
		}
		if queue.Blocked {
			return nil, backend.ErrBlocked
		}
		jobArg.ID = newID()
		jobArg.EnqueuedAt = now
		if jobArg.Version == nil {
			jobArg.Version = resource.NewVersion(0)
		}
		jobArg.Version.Inc()
		jobArg.QueueVersion = queue.Version
		jobArg.Status = resource.NewStatus(resource.StatusQueued)

		m.mu.Lock()
		m.jobs[jobArg.ID] = jobArg
		m.mu.Unlock()
	}
	return jobArgs, nil
}

func (m *Memory) DequeueJobs(ctx context.Context, limit int, opts *resource.JobListParams) (*resource.Jobs, error) {
	// fmt.Println("---\ndequeueJobs()")
	if opts == nil {
		opts = &resource.JobListParams{}
	}
	opts.Statuses = []*resource.Status{resource.NewStatus(resource.StatusQueued), resource.NewStatus(resource.StatusFailed)}

	jobs, err := m.ListJobs(ctx, limit, opts)
	if err != nil {
		return nil, err
	}

	now := internal.GetTimeProvider(ctx).Now().UTC()

	// filter out jobs with an unmet claim window or have failed too recently
	var filtered []*resource.Job
	for _, jb := range jobs.Jobs {
		queue, err := m.GetQueue(ctx, jb.Name, nil)
		if err != nil {
			return nil, err
		}
		if queue.Paused {
			continue
		}
		if jb.Data != nil && len(jb.Data.Claims) > 0 {
			match := jb.Data.Claims.Match(opts.Claims)
			expired := queue.ClaimExpired(jb, now)
			// fmt.Println("claim filter:", jb.ID, "match:", match, "expired:", expired)
			if !expired && !match {
				continue
			}
		}

		if len(jb.Results) > 0 {
			lastRes := jb.Results[len(jb.Results)-1]
			if lastRes.Status != nil && *lastRes.Status == resource.StatusFailed {
				if !lastRes.CompletedAt.IsZero() && now.After(lastRes.CompletedAt.Add(calculateBackoff(jb, queue))) {
					continue
				}
			}
		}

		filtered = append(filtered, jb)
	}
	jobs.Jobs = filtered

	if limit < len(jobs.Jobs) {
		jobs.Jobs = jobs.Jobs[:limit]
	}

	for _, jobData := range jobs.Jobs {
		status := resource.NewStatus(resource.StatusRunning)
		jobData.Version.Inc()
		jobData.Results = append(jobData.Results, &resource.JobResult{
			StartedAt: now,
			Status:    status,
		})
		jobData.Status = status
	}
	return jobs, nil
}

func (m *Memory) AckJobs(ctx context.Context, acks *resource.Acks) error {
	now := internal.GetTimeProvider(ctx).Now().UTC()
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, ack := range acks.Acks {
		jobData, ok := m.jobs[ack.JobID]
		if !ok {
			return backend.ErrNotFound
		}
		// if *jobData.Status != resource.StatusRunning {
		// 	// TODO return data about what state specifically caused this
		// 	return backend.ErrInvalidState
		// }

		jobData.Version.Inc()

		// fmt.Printf("ack %s: %#v\n", ack.ID, jobData)
		res := jobData.LastResult()
		res.CompletedAt = now
		if ack.Data != nil {
			res.Data = ack.Data
		}
		jobData.Status = ack.Status
	}
	// fmt.Println("---")
	// for k := range m.jobs {
	// 	fmt.Println(k, m.jobs[k].Status)
	// }
	return nil
}

// func (m *Memory) GetSetJobKeys(ctx context.Context, keys []string) (string, bool, error) {
// 	for _, key := range keys {
// 		_, ok := m.uniqueness[key]
// 		if ok {
// 			return "", true, nil
// 		}
// 	}
// 	for _, key := range keys {
// 		m.uniqueness[key] = true
// 	}
// 	return "", false, nil
// }

// func (m *Memory) DeleteJobKeys(ctx context.Context, keys []string) error {
// 	for _, key := range keys {
// 		delete(m.uniqueness, key)
// 	}
// 	return nil
// }

func (m *Memory) GetJobByID(ctx context.Context, id string, opts *resource.GetByIDOpts) (*resource.Job, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	jobData, ok := m.jobs[id]
	if !ok {
		return nil, backend.ErrNotFound
	}
	return jobData, nil
}

func (m *Memory) ListJobs(ctx context.Context, limit int, opts *resource.JobListParams) (*resource.Jobs, error) {
	if opts == nil {
		opts = &resource.JobListParams{}
	}
	res := &resource.Jobs{}
	m.mu.Lock()
	for _, jobData := range m.jobs {
		if statuses := opts.Statuses; len(statuses) > 0 && !jobData.HasStatus(statuses...) {
			continue
		}
		res.Jobs = append(res.Jobs, jobData)
	}
	m.mu.Unlock()

	if len(res.Jobs) > limit {
		res.Jobs = res.Jobs[:limit]
	}
	return res, nil
}

func (m *Memory) Stats(ctx context.Context, queue string) (*resource.Stats, error) {

	return nil, nil
}

func (m *Memory) GetJobUniqueArgs(ctx context.Context, keys []string) ([]string, bool, error) {
	return nil, false, nil
}

func (m *Memory) SetJobUniqueArgs(ctx context.Context, ids, keys []string) error {
	return nil
}

func (m *Memory) DeleteJobUniqueArgs(ctx context.Context, ids, keys []string) error {
	return nil
}

func (m *Memory) DeleteQueues(ctx context.Context, queues []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, queue := range queues {
		delete(m.queues, queue)
	}
	return nil
}

func (m *Memory) PauseQueues(ctx context.Context, queues []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	// TODO bump versions for these
	for _, queue := range queues {
		stored, ok := m.queues[queue]
		if !ok {
			return backend.ErrNotFound
		}
		stored.Paused = true
	}
	return nil
}

func (m *Memory) UnpauseQueues(ctx context.Context, queues []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, queue := range queues {
		stored, ok := m.queues[queue]
		if !ok {
			return backend.ErrNotFound
		}
		stored.Paused = false
		stored.Unpaused = true
	}
	return nil
}

func (m *Memory) BlockQueues(ctx context.Context, queues []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, queue := range queues {
		stored, ok := m.queues[queue]
		if !ok {
			return backend.ErrNotFound
		}
		stored.Blocked = true
	}
	return nil
}

func (m *Memory) UnblockQueues(ctx context.Context, queues []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, queue := range queues {
		stored, ok := m.queues[queue]
		if !ok {
			return backend.ErrNotFound
		}
		stored.Blocked = false
	}
	return nil
}

func valIn(val string, vals []string) bool {
	for _, v := range vals {
		if val == v {
			return true
		}
	}
	return false
}

func newID() string {
	return uuid.NewV4().String()
}

func calculateBackoff(jb *resource.Job, queue *resource.Queue) time.Duration {
	initial := time.Duration(queue.BackoffInitial)
	if initial == 0 {
		initial = time.Second
	}
	max := time.Duration(queue.BackoffMax)
	factor := queue.BackoffFactor
	if factor == 0 {
		factor = 2.0
	}
	attempt := jb.Attempt

	res := time.Duration(initial * time.Duration(math.Pow(float64(attempt), float64(factor))))
	if max > 0 && res > max {
		return max
	}
	return res
}
