package handler

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	apiv1 "github.com/jeffrom/job-manager/pkg/api/v1"
)

var internalServerErrorProto = &apiv1.GenericError{
	Message: "internal server error",
	Kind:    "internal",
}

var notFoundErrorProto = &apiv1.GenericError{
	Message: "not found",
	Kind:    "not_found",
}

type Error struct {
	status   int
	msg      string
	origErr  error
	resource string
	protoMsg proto.Message
}

func NewError(msg string) *Error {
	return &Error{
		msg:    msg,
		status: 500,
	}
}

func NewInternalServerError(err error) *Error {
	e := NewError("handler: internal server error")
	e.status = 500
	e.origErr = err
	e.protoMsg = internalServerErrorProto
	return e
}

func NewNotFoundError(resource string) *Error {
	msg := "not found"
	if resource != "" {
		msg = resource + msg
	}
	e := NewError(msg)
	e.status = 404
	e.resource = resource
	e.protoMsg = notFoundErrorProto
	return e
}

func (e *Error) Unwrap() error {
	return e.origErr
}

func (e *Error) Error() string {
	if e.origErr != nil {
		return fmt.Sprintf("%s: %v", e.msg, e.origErr)
	}
	return e.msg
}

func (e *Error) Status() int { return e.status }
