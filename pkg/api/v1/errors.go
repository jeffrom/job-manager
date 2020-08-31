package v1

import (
	"fmt"
	"net/http"

	"google.golang.org/protobuf/proto"
)

type NotFoundError struct {
	*NotFoundErrorResponse
}

func (e *NotFoundError) Is(other error) bool {
	_, ok := other.(*NotFoundError)
	return ok
}

func (e *NotFoundError) Error() string {
	if e.NotFoundErrorResponse == nil {
		return "api/v1: not found"
	}
	return fmt.Sprintf("api/v1: %q not found", e.NotFoundErrorResponse.Resource)
}

func (e *NotFoundError) Status() int { return http.StatusNotFound }

func (e *NotFoundError) Message() proto.Message { return e.NotFoundErrorResponse }

func NewNotFoundError(resource string) *NotFoundError {
	return &NotFoundError{
		NotFoundErrorResponse: &NotFoundErrorResponse{Resource: resource},
	}
}

func NewNotFoundErrorProto(message *NotFoundErrorResponse) *NotFoundError {
	return &NotFoundError{NotFoundErrorResponse: message}
}

// func NewValidationError(message *ValidationErrorResponse)
