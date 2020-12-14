package schema

import (
	"net/http"

	"github.com/hashicorp/go-multierror"
	"github.com/qri-io/jsonschema"

	"github.com/jeffrom/job-manager/mjob/resource"
)

type ValidationError struct {
	err    *multierror.Error
	Errors []jsonschema.KeyError
}

func (ve *ValidationError) Status() int {
	if ve != nil && ve.err != nil && len(ve.err.Errors) > 0 {
		return http.StatusBadRequest
	}
	return http.StatusOK
}

func (ve *ValidationError) Error() string {
	return ve.err.Error()
}

// func (ve *ValidationError) Message() proto.Message {
// 	msg := &apiv1.ValidationErrorResponse{
// 		Errs: []*apiv1.ValidationErrorArg{},
// 	}
// 	for _, verr := range ve.Errors {
// 		val, err := structpb.NewValue(verr.InvalidValue)
// 		if err != nil {
// 			return nil
// 		}
// 		msg.Errs = append(msg.Errs, &apiv1.ValidationErrorArg{
// 			Path:    verr.PropertyPath,
// 			Value:   val,
// 			Message: verr.Message,
// 		})
// 	}
// 	return msg
// }

func (ve *ValidationError) KeyErrors() []jsonschema.KeyError {
	return ve.Errors
}

func ErrorFromKeyErrors(resourceName, resourceID, reason string, errs []jsonschema.KeyError) *resource.Error {
	verrs := make([]*resource.ValidationError, len(errs))
	for i, serr := range errs {
		verrs[i] = &resource.ValidationError{
			Message: serr.Message,
			Path:    serr.PropertyPath,
			Value:   serr.InvalidValue,
		}
	}
	return resource.NewValidationError(resourceName, resourceID, reason, verrs)
}

func NewValidationErrorKeyErrs(errs []jsonschema.KeyError) *ValidationError {
	ve := &ValidationError{Errors: errs, err: &multierror.Error{}}
	for _, verr := range errs {
		ve.err.Errors = append(ve.err.Errors, verr)
	}

	return ve
}
