package v1

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/proto"
)

type Error struct {
	status     int
	msg        string
	kind       string
	resource   string
	resourceID string
	reason     string
	origErr    error
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

func NewConflictError(resource, resourceID, reason string) *Error {
	msg := "conflict"
	if resource != "" {
		msg = resource + " " + msg
	}
	return &Error{
		status:     409,
		msg:        msg,
		kind:       "conflict",
		resource:   resource,
		resourceID: resourceID,
		reason:     reason,
	}
}

func NewUnprocessableEntityError(resource string) *Error {
	return &Error{
		status:   422,
		msg:      "unprocessable entity",
		kind:     "unprocessable_entity",
		resource: resource,
	}
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

func ErrorFromProto(msg proto.Message) (*Error, error) {
	return nil, nil
}

func (e *Error) Error() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("api/v1 %s: %s", e.kind, e.msg))
	if e.reason != "" {
		b.WriteString(". ")
		b.WriteString(e.reason)
	}
	if e.resource != "" {
		b.WriteString(" (resource: ")
		b.WriteString(e.resource)
	}
	if e.resourceID != "" {
		b.WriteString(", id: ")
		b.WriteString(e.resourceID)
	}
	if e.resource != "" || e.resourceID != "" {
		b.WriteString(")")
	}

	if orig := e.origErr; orig != nil {
		b.WriteString(fmt.Sprintf(" (original error: %v)", orig))
	}

	return b.String()
}

func (e *Error) Unwrap() error { return e.origErr }

func (e *Error) Is(other error) bool {
	if e == nil || other == nil {
		return e == other
	}

	otherGeneric, ok := other.(*Error)
	if !ok {
		return false
	}

	if e.kind != otherGeneric.kind {
		return false
	}
	if e.resource != otherGeneric.resource {
		return false
	}
	if e.resourceID != otherGeneric.resourceID {
		return false
	}

	return true
}

func (e *Error) Status() int { return e.status }

func (e *Error) Message() proto.Message {
	return &GenericError{
		Kind:       e.kind,
		Message:    e.msg,
		Resource:   e.resource,
		ResourceId: e.resourceID,
		Reason:     e.reason,
	}
}
