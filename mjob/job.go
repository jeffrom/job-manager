package mjob

import "github.com/jeffrom/job-manager/mjob/resource"

// RunSucceeded returns a successful job result with optional data. It is
// a convenience function intended to be used as the return value to
// implementations of Runner.
func RunSucceeded(data interface{}) (*Result, error) {
	res := &resource.JobResult{
		Status: resource.NewStatus(resource.StatusComplete),
		Data:   data,
	}

	return res, nil
}

func RunInvalid(err error, data interface{}) (*Result, error) {
	panic("not implemented")
}

// RunFailed returns a failed job result with optional error and data. It is
// a convenience function intended to be used as the return value to
// implementations of Runner.
func RunFailed(err error, data interface{}) (*Result, error) {
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
