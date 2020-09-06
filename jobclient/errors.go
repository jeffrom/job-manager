package jobclient

import (
	"errors"

	apiv1 "github.com/jeffrom/job-manager/pkg/api/v1"
	"github.com/jeffrom/job-manager/pkg/schema"
	"github.com/qri-io/jsonschema"
)

var ErrInternal = errors.New("jobclient: internal server error")

type APIError struct {
	*apiv1.GenericError
}

func (e *APIError) Error() string {
	return apiv1.ErrorMessage(e.GenericError)
}

func (e *APIError) Is(other error) bool {
	if e == nil || other == nil {
		return e == other
	}

	otherGeneric, ok := other.(*APIError)
	if !ok {
		return false
	}

	if e.Kind != otherGeneric.Kind {
		return false
	}
	if e.Resource != otherGeneric.Resource {
		return false
	}
	if e.ResourceId != otherGeneric.ResourceId {
		return false
	}

	return true
}

func newGenericErrorFromMessage(message *apiv1.GenericError) *APIError {
	return &APIError{GenericError: message}
}

func IsNotFound(err error) bool {
	gerr := &APIError{}
	if errors.As(err, &gerr) {
		return gerr.Kind == "not_found"
	}
	return false
}

func newSchemaValidationErrorProto(resp *apiv1.ValidationErrorResponse) *schema.ValidationError {
	ve := &schema.ValidationError{}
	for _, verr := range resp.Errs {
		ve.Errors = append(ve.Errors, jsonschema.KeyError{
			PropertyPath: verr.Path,
			InvalidValue: verr.Value.AsInterface(),
			Message:      verr.Message,
		})
	}
	return ve
}
