package mjob

import "github.com/jeffrom/job-manager/mjob/resource"

// RunSuccess returns a successful job result with optional data. It is
// intended to be used as the return value to implementations of Runner.
func RunSuccess(data interface{}) (*resource.JobResult, error) {
	res := &resource.JobResult{
		Status: resource.NewStatus(resource.StatusComplete),
		Data:   data,
	}

	return res, nil
}

func RunInvalid(err error, data interface{}) (*resource.JobResult, error) {
	panic("not implemented")
}

// RunFail returns a failed job result with optional error and data. It is
// intended to be used as the return value to implementations of Runner.
func RunFail(err error, data interface{}) (*resource.JobResult, error) {
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
