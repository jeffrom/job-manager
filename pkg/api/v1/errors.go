package v1

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

type Error struct {
	status   int
	msg      string
	kind     string
	resource string
	origErr  error
}

func NewError(msg string) *Error {
	return &Error{
		msg:    msg,
		kind:   "internal",
		status: 500,
	}
}

func NewInternalServerError(err error) *Error {
	e := NewError("handler: internal server error")
	e.status = 500
	e.kind = "internal"
	e.origErr = err
	return e
}

func NewConflictError(resource string) *Error {
	msg := "conflict"
	if resource != "" {
		msg = resource + " " + msg
	}
	e := NewError(msg)
	e.kind = "conflict"
	e.status = 409
	return e
}

func NewNotFoundError(resource string) *Error {
	msg := "not found"
	if resource != "" {
		msg = resource + " " + msg
	}
	e := NewError(msg)
	e.kind = "not_found"
	e.resource = resource
	e.status = 404
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

func (e *Error) Message() proto.Message {
	return &GenericError{
		Kind:     e.kind,
		Message:  e.msg,
		Resource: e.resource,
	}
}
