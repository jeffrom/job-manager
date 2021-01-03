package client

import (
	"errors"

	"github.com/qri-io/jsonschema"

	apiv1 "github.com/jeffrom/job-manager/mjob/api/v1"
	"github.com/jeffrom/job-manager/mjob/resource"
	"github.com/jeffrom/job-manager/mjob/schema"
)

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

func newResourceErrorFromMessage(message *apiv1.GenericError) *resource.Error {
	var code int = 500
	switch message.Kind {
	case "not_found":
		code = 404
	case "internal":
		code = 500
	case "conflict":
		code = 409
	case "unprocessable_entity":
		code = 422
	case "invalid":
		code = 400
	}

	msg := message.Message
	if msg == "" {
		msg = "internal server error"
	}
	return &resource.Error{
		Status:     code,
		Kind:       message.Kind,
		Message:    msg,
		Resource:   message.Resource,
		ResourceID: message.ResourceId,
		Reason:     message.Reason,
		Invalid:    newResourceValidationErrorsProto(message.Invalid),
	}
}

func IsNotFound(err error) bool {
	gerr := &resource.Error{}
	if errors.As(err, &gerr) {
		return gerr.Kind == "not_found"
	}
	return false
}

func newResourceValidationErrorsProto(resp []*apiv1.ValidationError) []*resource.ValidationError {
	ve := make([]*resource.ValidationError, len(resp))
	for i, respErr := range resp {
		ve[i] = &resource.ValidationError{
			Path:    respErr.Path,
			Message: respErr.Message,
			Value:   respErr.Value.AsInterface(),
		}
	}
	return ve
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
