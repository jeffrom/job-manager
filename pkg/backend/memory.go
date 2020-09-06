package backend

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"

	"github.com/jeffrom/job-manager/pkg/label"
	jobv1 "github.com/jeffrom/job-manager/pkg/resource/job/v1"
)

// Memory is an in-memory backend intended for testing.
type Memory struct {
	configs    map[string]*jobv1.Queue
	jobs       map[string]*jobv1.Job
	uniqueness map[string]bool
}

func NewMemory() *Memory {
	return &Memory{
		configs:    make(map[string]*jobv1.Queue),
		jobs:       make(map[string]*jobv1.Job),
		uniqueness: make(map[string]bool),
	}
}

func (m *Memory) GetQueue(ctx context.Context, name string) (*jobv1.Queue, error) {
	if cfg, ok := m.configs[name]; ok {
		return cfg, nil
	}
	return nil, ErrNotFound
}

func (m *Memory) SaveQueue(ctx context.Context, queue *jobv1.Queue) error {
	// if it already exists and no version was supplied, or if version was
	// supplied but they don't match, return conflict
	prev, ok := m.configs[queue.Id]
	// fmt.Printf("prev: %+v, found: %v\n", prev, ok)
	if ok {
		if queue.V == 0 || queue.V != prev.V {
			return &VersionConflictError{
				Resource:   "queue",
				ResourceID: queue.Id,
				Prev:       strconv.FormatInt(int64(prev.V), 10),
				Curr:       strconv.FormatInt(int64(queue.V), 10),
			}
		}
	}
	// fmt.Printf("---\nprev: %s\ncurr: %s\nequal: %v\n\n", prev.String(), queue.String(), proto.Equal(queue, prev))
	if prev == nil || !queuesEqual(queue, prev) {
		queue.V++
	}
	m.configs[queue.Id] = queue
	return nil
}

func (m *Memory) ListQueues(ctx context.Context, opts *jobv1.QueueListParams) (*jobv1.Queues, error) {
	if opts == nil {
		opts = &jobv1.QueueListParams{}
	}
	sels, err := label.ParseSelectorStringArray(opts.Selectors)
	if err != nil {
		return nil, err
	}

	queues := &jobv1.Queues{}
	for _, queue := range m.configs {
		if m.filterQueue(queue, opts.Names, sels) {
			continue
		}
		queues.Queues = append(queues.Queues, queue)
	}
	return queues, nil
}

func (m *Memory) filterQueue(queue *jobv1.Queue, names []string, sels *label.Selectors) bool {
	if len(names) > 0 && !valIn(queue.Id, names) {
		return true
	}
	return !sels.Match(queue.Labels)
}

func (m *Memory) EnqueueJobs(ctx context.Context, jobArgs *jobv1.Jobs) error {
	for _, jobArg := range jobArgs.Jobs {
		m.jobs[jobArg.Id] = jobArg
	}
	return nil
}

func (m *Memory) DequeueJobs(ctx context.Context, num int, opts *jobv1.JobListParams) (*jobv1.Jobs, error) {
	if opts == nil {
		opts = &jobv1.JobListParams{}
	}
	opts.Statuses = []jobv1.Status{jobv1.StatusQueued, jobv1.StatusFailed}

	jobs, err := m.ListJobs(ctx, opts)
	if err != nil {
		return nil, err
	}
	if num < len(jobs.Jobs) {
		jobs.Jobs = jobs.Jobs[:num]
	}

	for _, jobData := range jobs.Jobs {
		jobData.Status = jobv1.StatusRunning
	}
	return jobs, nil
}

func (m *Memory) AckJobs(ctx context.Context, results *jobv1.Acks) error {
	for _, res := range results.Acks {
		jobData, ok := m.jobs[res.Id]
		if !ok {
			return ErrNotFound
		}
		if res.Data != nil {
			jobData.Results = []*jobv1.Result{
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

func (m *Memory) GetJobByID(ctx context.Context, id string) (*jobv1.Job, error) {
	jobData, ok := m.jobs[id]
	if !ok {
		return nil, ErrNotFound
	}
	return jobData, nil
}

func (m *Memory) ListJobs(ctx context.Context, opts *jobv1.JobListParams) (*jobv1.Jobs, error) {
	if opts == nil {
		opts = &jobv1.JobListParams{}
	}
	res := &jobv1.Jobs{}
	for _, jobData := range m.jobs {
		if statuses := opts.Statuses; len(statuses) > 0 && !jobv1.HasStatus(jobData, statuses) {
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

func queuesEqual(a, b *jobv1.Queue) bool {
	if a == nil || b == nil {
		return a == b
	}
	ac := &*a
	bc := &*b
	ac.CreatedAt = nil
	bc.CreatedAt = nil
	abuf, err := json.Marshal(ac)
	if err != nil {
		return false
	}
	bbuf, err := json.Marshal(bc)
	if err != nil {
		return false
	}
	return bytes.Equal(abuf, bbuf)
}

func labelsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
