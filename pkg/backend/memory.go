package backend

import (
	"context"

	"github.com/jeffrom/job-manager/pkg/job"
)

// Memory is an in-memory backend intended for testing.
type Memory struct {
	configs map[string]*job.Queue
	jobs    map[string]*job.Job
}

func NewMemory() *Memory {
	return &Memory{
		configs: make(map[string]*job.Queue),
		jobs:    make(map[string]*job.Job),
	}
}

func (m *Memory) GetQueue(ctx context.Context, name string) (*job.Queue, error) {
	if cfg, ok := m.configs[name]; ok {
		return cfg, nil
	}
	return nil, ErrNotFound
}

func (m *Memory) SaveQueue(ctx context.Context, queue *job.Queue) error {
	m.configs[queue.Id] = queue
	return nil
}

func (m *Memory) ListQueues(ctx context.Context, opts *job.QueueListParams) (*job.Queues, error) {
	jobs := &job.Queues{}
	for _, job := range m.configs {
		jobs.Queues = append(jobs.Queues, job)
	}
	return jobs, nil
}

func (m *Memory) EnqueueJobs(ctx context.Context, jobArgs *job.Jobs) error {
	for _, jobArg := range jobArgs.Jobs {
		m.jobs[jobArg.Id] = jobArg
	}
	return nil
}

func (m *Memory) DequeueJobs(ctx context.Context, num int, opts *job.JobListParams) (*job.Jobs, error) {
	if opts == nil {
		opts = &job.JobListParams{}
	}
	opts.Statuses = []job.Status{job.StatusQueued, job.StatusFailed}

	jobs, err := m.ListJobs(ctx, opts)
	if err != nil {
		return nil, err
	}
	if num < len(jobs.Jobs) {
		jobs.Jobs = jobs.Jobs[:num]
	}

	for _, jobData := range jobs.Jobs {
		jobData.Status = job.StatusRunning
	}
	return jobs, nil
}

func (m *Memory) AckJobs(ctx context.Context, results *job.Results) error {
	for _, res := range results.Results {
		job, ok := m.jobs[res.Id]
		if !ok {
			return ErrNotFound
		}
		if len(res.Data) > 0 {
			job.ResultData = res.Data
		}
		job.Status = res.Status
	}
	// fmt.Println("---")
	// for k := range m.jobs {
	// 	fmt.Println(k, m.jobs[k].Status)
	// }
	return nil
}

func (m *Memory) GetJobByID(ctx context.Context, id string) (*job.Job, error) {
	jobData, ok := m.jobs[id]
	if !ok {
		return nil, ErrNotFound
	}
	return jobData, nil
}

func (m *Memory) ListJobs(ctx context.Context, opts *job.JobListParams) (*job.Jobs, error) {
	if opts == nil {
		opts = &job.JobListParams{}
	}
	res := &job.Jobs{}
	for _, jobData := range m.jobs {
		if statuses := opts.Statuses; len(statuses) > 0 && !job.HasStatus(jobData, statuses) {
			continue
		}
		res.Jobs = append(res.Jobs, jobData)
	}
	return res, nil
}
