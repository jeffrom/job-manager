package schema

import (
	"net/http"

	"github.com/hashicorp/go-multierror"
	"github.com/qri-io/jsonschema"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"

	apiv1 "github.com/jeffrom/job-manager/pkg/api/v1"
)

type ValidationError struct {
	err   *multierror.Error
	verrs []jsonschema.KeyError
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

func (ve *ValidationError) Message() proto.Message {
	msg := &apiv1.ValidationErrorResponse{
		Errs: []*apiv1.ValidationErrorArg{},
	}
	for _, verr := range ve.verrs {
		val, err := structpb.NewValue(verr.InvalidValue)
		if err != nil {
			return nil
		}
		msg.Errs = append(msg.Errs, &apiv1.ValidationErrorArg{
			Path:    verr.PropertyPath,
			Value:   val,
			Message: verr.Message,
		})
	}
	return msg
}

func (ve *ValidationError) KeyErrors() []jsonschema.KeyError {
	return ve.verrs
}

func NewValidationErrorKeyErrs(errs []jsonschema.KeyError) *ValidationError {
	ve := &ValidationError{verrs: errs, err: &multierror.Error{}}
	for _, verr := range errs {
		ve.err.Errors = append(ve.err.Errors, verr)
	}

	return ve
}

func NewValidationErrorProto(resp *apiv1.ValidationErrorResponse) *ValidationError {
	ve := &ValidationError{}
	for _, verr := range resp.Errs {
		ve.verrs = append(ve.verrs, jsonschema.KeyError{
			PropertyPath: verr.Path,
			InvalidValue: verr.Value.AsInterface(),
			Message:      verr.Message,
		})
	}
	return ve
}
