package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/jeffrom/job-manager/mjob/resource"
	"github.com/jeffrom/job-manager/mjob/schema"
	"github.com/jeffrom/job-manager/pkg/backend"
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
	case errors.Is(err, backend.ErrInvalidState):
		return resource.NewUnprocessableEntityError(resourceName, resourceID, "job state does not allow this operation")
	case errors.Is(err, backend.ErrNotFound):
		return resource.NewNotFoundError(resourceName, resourceID, "")
	case errors.As(err, &cerr):
		return newVersionConflictError(cerr)
	}
	return err
}

func handleSchemaErrors(err error, resourceName, resourceID, reason string) error {
	verr := &schema.ValidationError{}
	if errors.As(err, &verr) {
		// fmt.Printf("handler: %#v\n", verr)
		return schema.ErrorFromKeyErrors(resourceName, resourceID, reason, verr.Errors)
	}
	return err
}
