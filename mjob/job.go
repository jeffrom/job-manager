package mjob

import "github.com/jeffrom/job-manager/mjob/resource"

type Job struct {
	*resource.Job
}

func (jb *Job) RunSuccess(data interface{}) (*resource.JobResult, error) {
	res := &resource.JobResult{
		Status: resource.NewStatus(resource.StatusComplete),
		Data:   data,
	}

	return res, nil
}

func (jb *Job) RunInvalid(err error, data interface{}) (*resource.JobResult, error) {
	panic("not implemented")
}

func (jb *Job) RunFail(err error, data interface{}) (*resource.JobResult, error) {
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	res := &resource.JobResult{
		Error:  errStr,
		Status: resource.NewStatus(resource.StatusFailed),
		Data:   data,
	}

	return res, nil
}
