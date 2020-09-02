package v1

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

var internalServerErrorProto = &GenericError{
	Message: "internal server error",
}

var notFoundErrorProto = &GenericError{
	Message: "not found",
}

type Error struct {
	status   int
	msg      string
	origErr  error
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

func (e *Error) Message() proto.Message { return e.protoMsg }
