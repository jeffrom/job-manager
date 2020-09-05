package jobclient

import (
	"errors"
	"fmt"

	apiv1 "github.com/jeffrom/job-manager/pkg/api/v1"
)

var ErrInternal = errors.New("jobclient: internal server error")

type GenericError struct {
	*apiv1.GenericError
}

func (e *GenericError) Error() string {
	if e.Resource != "" {
		return fmt.Sprintf("api/v1: %s %s", e.Resource, e.Message)
	}
	return fmt.Sprintf("api/v1: %s", e.Message)
}

func (e *GenericError) Is(other error) bool {
	if e == nil || other == nil {
		return e == other
	}

	otherGeneric, ok := other.(*GenericError)
	if !ok {
		return false
	}

	if e.Kind != otherGeneric.Kind {
		return false
	}
	return e.Resource == otherGeneric.Resource
}

func newGenericErrorFromMessage(message *apiv1.GenericError) *GenericError {
	return &GenericError{GenericError: message}
}
