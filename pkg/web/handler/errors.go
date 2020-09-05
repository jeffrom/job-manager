package handler

import (
	"errors"
	"fmt"
	"net/http"

	apiv1 "github.com/jeffrom/job-manager/pkg/api/v1"
	"github.com/jeffrom/job-manager/pkg/backend"
)

func newVersionConflictError(e *backend.VersionConflictError) *apiv1.Error {
	reason := fmt.Sprintf("a conflicting version was provided (existing: %q, current: %q)", e.Prev, e.Curr)
	return apiv1.NewConflictError(e.Resource, e.ResourceID, reason)
}

func NotFound(w http.ResponseWriter, r *http.Request) error {
	return apiv1.ErrGenericNotFound
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) error {
	return apiv1.ErrMethodNotAllowed
}

func handleBackendErrors(err error, resource, resourceID string) error {
	cerr := &backend.VersionConflictError{}
	switch {
	case errors.Is(err, backend.ErrNotFound):
		return apiv1.NewNotFoundError(resource, resourceID, "")
	case errors.As(err, &cerr):
		return newVersionConflictError(cerr)
	}
	return err
}
