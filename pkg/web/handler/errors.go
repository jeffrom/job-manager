package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/jeffrom/job-manager/pkg/backend"
	"github.com/jeffrom/job-manager/pkg/resource"
)

func newVersionConflictError(e *backend.VersionConflictError) *resource.Error {
	reason := fmt.Sprintf("a conflicting version was provided (existing: %q, current: %q)", e.Prev, e.Curr)
	return resource.NewConflictError(e.Resource, e.ResourceID, reason)
}

func NotFound(w http.ResponseWriter, r *http.Request) error {
	return resource.ErrGenericNotFound
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) error {
	return resource.ErrMethodNotAllowed
}

func handleBackendErrors(err error, resourceName, resourceID string) error {
	cerr := &backend.VersionConflictError{}
	switch {
	case errors.Is(err, backend.ErrNotFound):
		return resource.NewNotFoundError(resourceName, resourceID, "")
	case errors.As(err, &cerr):
		return newVersionConflictError(cerr)
	}
	return err
}
