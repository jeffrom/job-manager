package jobclient

import (
	"errors"
	"fmt"

	apiv1 "github.com/jeffrom/job-manager/pkg/api/v1"
)

var ErrInternal = errors.New("jobclient: internal server error")

type NotFoundError struct {
	*apiv1.GenericError
	// Resource string `json:"resource"`
}

func (e *NotFoundError) Is(other error) bool {
	_, ok := other.(*NotFoundError)
	return ok
}

func (e *NotFoundError) Error() string {
	if e.GenericError == nil {
		return "api/v1: not found"
	}
	return fmt.Sprintf("api/v1: %q not found", e.GenericError.Resource)
}

func newNotFoundErrorProto(message *apiv1.GenericError) *NotFoundError {
	return &NotFoundError{GenericError: message}
}
