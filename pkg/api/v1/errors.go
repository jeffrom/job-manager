package v1

import (
	"fmt"
	"strings"

	"github.com/jeffrom/job-manager/pkg/resource"
	"google.golang.org/protobuf/proto"
)

func ErrorMessage(e *GenericError) string {
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
	if e.ResourceId != "" {
		b.WriteString(", id: ")
		b.WriteString(e.ResourceId)
	}
	if e.Resource != "" || e.ResourceId != "" {
		b.WriteString(")")
	}

	// if orig := e.origErr; orig != nil {
	// 	b.WriteString(fmt.Sprintf(" (original error: %v)", orig))
	// }

	return b.String()
}

type protoError struct {
	rerr *resource.Error
}

func (e *protoError) Error() string {
	return e.rerr.Error()
}

func (e *protoError) Message() proto.Message {
	return &GenericError{
		Message:    e.rerr.Message,
		Kind:       e.rerr.Kind,
		Resource:   e.rerr.Resource,
		ResourceId: e.rerr.ResourceID,
		Reason:     e.rerr.Reason,
	}
}

func ErrorProto(rerr *resource.Error) *protoError {
	return &protoError{
		rerr: rerr,
	}
}
