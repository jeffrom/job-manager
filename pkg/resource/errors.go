package resource

import (
	"fmt"
	"net/http"
	"strings"
)

var ErrMethodNotAllowed = &Error{
	Status:  http.StatusMethodNotAllowed,
	Message: "method not allowed",
	Kind:    "method_not_allowed",
}

var ErrGenericNotFound = &Error{
	Status:  http.StatusNotFound,
	Message: "not found",
	Kind:    "not_found",
}

type Error struct {
	Status     int                `json:"status,omitempty"`
	Message    string             `json:"message"`
	Kind       string             `json:"kind"`
	Resource   string             `json:"resource,omitempty"`
	ResourceID string             `json:"resource_id,omitempty"`
	Reason     string             `json:"reason,omitempty"`
	Orig       error              `json:"orig,omitempty"`
	Invalid    []*ValidationError `json:"invalid,omitempty"`
}

func NewInternalServerError(err error) *Error {
	return &Error{
		Kind:    "internal",
		Message: "internal server error",
		Orig:    err,
	}
}

func NewConflictError(resource, resourceID, reason string) *Error {
	msg := "conflict"
	if resource != "" {
		msg = resource + " " + msg
	}
	return &Error{
		Status:     409,
		Message:    msg,
		Kind:       "conflict",
		Resource:   resource,
		ResourceID: resourceID,
		Reason:     reason,
	}
}

func NewUnprocessableEntityError(resource, resourceID, reason string) *Error {
	return &Error{
		Status:     422,
		Message:    "unprocessable entity",
		Kind:       "unprocessable_entity",
		Resource:   resource,
		ResourceID: resourceID,
		Reason:     reason,
	}
}

func NewNotFoundError(resource, resourceID, reason string) *Error {
	return &Error{
		Status:     404,
		Message:    "not found",
		Kind:       "not_found",
		Resource:   resource,
		ResourceID: resourceID,
		Reason:     reason,
	}
}

func (e *Error) Error() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("api/v1: %s", e.Message))
	if e.Reason != "" {
		b.WriteString(": ")
		b.WriteString(e.Reason)
	}
	if e.Resource != "" {
		b.WriteString(" (resource: ")
		b.WriteString(e.Resource)
	}
	if e.ResourceID != "" {
		b.WriteString(", id: ")
		b.WriteString(e.ResourceID)
	}
	if e.Resource != "" || e.ResourceID != "" {
		b.WriteString(")")
	}

	if orig := e.Orig; orig != nil {
		b.WriteString(fmt.Sprintf(" (original error: %v)", orig))
	}

	return b.String()
}

func (e *Error) Unwrap() error { return e.Orig }

func (e *Error) Is(iother error) bool {
	if e == nil || iother == nil {
		return e == iother
	}

	other, ok := iother.(*Error)
	if !ok {
		return false
	}

	if e.Kind != other.Kind {
		return false
	}
	if e.Resource != other.Resource {
		return false
	}
	if e.ResourceID != other.ResourceID {
		return false
	}
	if e.Orig != other.Orig {
		return false
	}

	return true
}

func (e *Error) GetStatus() int { return e.Status }

type ValidationError struct {
	Message  string      `json:"message"`
	Path     string      `json:"path,omitempty"`
	ValueRaw string      `json:"value_raw,omitempty"`
	Value    interface{} `json:"value,omitempty"`
}
