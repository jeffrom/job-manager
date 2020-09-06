package backend

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("backend: not found")

type VersionConflictError struct {
	Resource   string
	ResourceID string
	Prev       string
	Curr       string
}

func (e *VersionConflictError) Error() string {
	return fmt.Sprintf("backend: version conflict (resource: %q prev: %q, curr: %q)", e.Resource, e.Prev, e.Curr)
}
