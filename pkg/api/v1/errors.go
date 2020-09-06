package v1

import (
	"fmt"
	"strings"
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
