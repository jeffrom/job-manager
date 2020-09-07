package backend

import (
	"context"
	"fmt"

	"github.com/jeffrom/job-manager/pkg/label"
	"github.com/jeffrom/job-manager/pkg/resource"
)

// Memory is an in-memory backend intended for testing.
type Memory struct {
	// mu         sync.Mutex
	queues     map[string]*resource.Queue
	jobs       map[string]*resource.Job
	uniqueness map[string]bool
}

func NewMemory() *Memory {
	return &Memory{
		queues:     make(map[string]*resource.Queue),
		jobs:       make(map[string]*resource.Job),
		uniqueness: make(map[string]bool),
	}
}

func (m *Memory) GetQueue(ctx context.Context, name string) (*resource.Queue, error) {
	if cfg, ok := m.queues[name]; ok {
		return cfg, nil
	}
	return nil, ErrNotFound
}

func (m *Memory) SaveQueue(ctx context.Context, queue *resource.Queue) error {
	// if it already exists and no version was supplied, or if version was
	// supplied but they don't match, return conflict
	prev, ok := m.queues[queue.ID]
	// fmt.Printf("prev: %+v, found: %v\n", prev, ok)
	if ok {
		if queue.Version.Raw() == 0 || queue.Version.Raw() != prev.Version.Raw() {
			return &VersionConflictError{
				Resource:   "queue",
				ResourceID: queue.ID,
				Prev:       prev.Version.String(),
				Curr:       queue.Version.String(),
			}
		}
	}
	// fmt.Printf("---\nprev: %s\ncurr: %s\nequal: %v\n\n", prev.String(), queue.String(), proto.Equal(queue, prev))
	if prev == nil || !queue.Equals(prev) {
		queue.Version.Inc()
	}
	m.queues[queue.ID] = queue
	return nil
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
	if len(names) > 0 && !valIn(queue.ID, names) {
		return true
	}
	return !sels.Match(queue.Labels)
}

func (m *Memory) EnqueueJobs(ctx context.Context, jobArgs *resource.Jobs) error {
	for _, jobArg := range jobArgs.Jobs {
		m.jobs[jobArg.ID] = jobArg
	}
	return nil
}

func (m *Memory) DequeueJobs(ctx context.Context, num int, opts *resource.JobListParams) (*resource.Jobs, error) {
	if opts == nil {
		opts = &resource.JobListParams{}
	}
	opts.Statuses = []resource.Status{resource.StatusQueued, resource.StatusFailed}

	jobs, err := m.ListJobs(ctx, opts)
	if err != nil {
		return nil, err
	}

	// filter out jobs with an unmet claim window
	for _, jb := range jobs.Jobs {
		fmt.Printf("%+v\n", jb)
	}

	if num < len(jobs.Jobs) {
		jobs.Jobs = jobs.Jobs[:num]
	}

	for _, jobData := range jobs.Jobs {
		jobData.Status = resource.StatusRunning
	}
	return jobs, nil
}

func (m *Memory) AckJobs(ctx context.Context, results *resource.Acks) error {
	for _, res := range results.Acks {
		jobData, ok := m.jobs[res.ID]
		if !ok {
			return ErrNotFound
		}
		if res.Data != nil {
			jobData.Results = []*resource.JobResult{
				{
					Data: res.Data,
				},
			}
		}
		jobData.Status = res.Status
	}
	// fmt.Println("---")
	// for k := range m.jobs {
	// 	fmt.Println(k, m.jobs[k].Status)
	// }
	return nil
}

func (m *Memory) GetSetJobKeys(ctx context.Context, keys []string) (bool, error) {
	for _, key := range keys {
		_, ok := m.uniqueness[key]
		m.uniqueness[key] = true
		if ok {
			return true, nil
		}
	}
	return false, nil
}

func (m *Memory) DeleteJobKeys(ctx context.Context, keys []string) error {
	for _, key := range keys {
		delete(m.uniqueness, key)
	}
	return nil
}

func (m *Memory) GetJobByID(ctx context.Context, id string) (*resource.Job, error) {
	jobData, ok := m.jobs[id]
	if !ok {
		return nil, ErrNotFound
	}
	return jobData, nil
}

func (m *Memory) ListJobs(ctx context.Context, opts *resource.JobListParams) (*resource.Jobs, error) {
	if opts == nil {
		opts = &resource.JobListParams{}
	}
	res := &resource.Jobs{}
	for _, jobData := range m.jobs {
		if statuses := opts.Statuses; len(statuses) > 0 && !jobData.HasStatus(statuses...) {
			continue
		}
		res.Jobs = append(res.Jobs, jobData)
	}
	return res, nil
}

func valIn(val string, vals []string) bool {
	for _, v := range vals {
		if val == v {
			return true
		}
	}
	return false
}
