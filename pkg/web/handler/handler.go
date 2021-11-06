// Package handler contains http handlers.
package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"google.golang.org/protobuf/proto"

	apiv1 "github.com/jeffrom/job-manager/mjob/api/v1"
	"github.com/jeffrom/job-manager/mjob/resource"
	"github.com/jeffrom/job-manager/pkg/backend"
	"github.com/jeffrom/job-manager/pkg/logger"
)

type httpError interface {
	error
	GetStatus() int
}

type protoError interface {
	Message() proto.Message
}

func Func(fn func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			reqLog := logger.RequestLogFromContext(r.Context())
			// fmt.Printf("req error: %T %+v\n", err, err)

			// set status depending on the error returned
			status := http.StatusInternalServerError
			if herr, ok := err.(httpError); ok {
				status = herr.GetStatus()
			} else if errors.Is(err, io.ErrUnexpectedEOF) {
				status = http.StatusBadRequest
			}

			reqLog.Str("err_type", fmt.Sprintf("%T", err))
			if status >= http.StatusInternalServerError {
				reqLog.Err(err)
			}

			// maybe translate errors into api errors
			// XXX this shouldn't be needed
			vcErr := &backend.VersionConflictError{}
			if errors.As(err, &vcErr) {
				err = newVersionConflictError(vcErr)
			}

			rerr := &resource.Error{}
			if errors.As(err, &rerr) {
				err = apiv1.ErrorProto(rerr)
			}

			// log the error
			log := logger.FromContext(r.Context())
			logRequestError(log, status, err)
			w.WriteHeader(status)

			if pr, ok := err.(protoError); ok {
				// fmt.Printf("handler: %T %+v\n", pr.Message(), pr.Message())
				if err := MarshalResponse(w, r, pr.Message()); err != nil {
					logger.FromContext(r.Context()).Error().Err(err).Msg("marshal response failed")
				}
			}
		}
	}
}

func UnmarshalBody(r *http.Request, v interface{}, required bool) error {
	log := logger.FromContext(r.Context())
	defer r.Body.Close()
	ct := r.Header.Get("content-type")
	var b []byte
	if r.Method == "GET" {
		ct = "application/x-www-form-urlencoded"
	} else {
		var rerr error
		b, rerr = io.ReadAll(r.Body)
		if rerr != nil {
			return rerr
		}
		if len(b) == 0 {
			if required {
				return io.ErrUnexpectedEOF
			}
			return nil
		}
	}

	var err error
	switch ct {
	case "", "application/json":
		err = json.Unmarshal(b, v)
	case "application/protobuf":
		err = proto.Unmarshal(b, v.(proto.Message))
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			return err
		}
		err = formDecoder.Decode(v, r.Form)
	default:
		panic("handler: unknown content-type: " + ct)
	}

	log.Debug().
		Str("content_type", ct).
		Str("type", fmt.Sprintf("%T", v)).
		Interface("params", v).
		Msg("params")
	return err
}

func MarshalResponse(w http.ResponseWriter, r *http.Request, v proto.Message) error {
	ct := r.Header.Get("content-type")
	var b []byte
	var err error
	switch ct {
	case "", "application/json":
		b, err = json.Marshal(v)
	case "application/protobuf":
		b, err = proto.Marshal(v)
	default:
		panic("handler: unknown content-type: " + ct)
	}

	log := logger.FromContext(r.Context())
	log.Debug().
		Str("type", fmt.Sprintf("%T", v)).
		Msg("response")
	if _, err := w.Write(b); err != nil {
		return err
	}
	return err
}

func logRequestError(log *logger.Logger, status int, err error) {
	ev := log.Warn()
	if status >= http.StatusInternalServerError {
		ev = log.Error()
	}

	ev.Err(err)

	ev.Msg("request error")
}

func readPaginationFromForm(form url.Values) (*apiv1.Pagination, error) {
	limitS := form.Get("page[limit]")
	var limit int64
	if limitS != "" {
		n, err := strconv.ParseInt(limitS, 10, 64)
		if err != nil {
			return nil, err
		}
		limit = n
	}

	lastID := form.Get("page[last_id]")
	if limit > 0 || lastID != "" {
		return &apiv1.Pagination{Limit: limit, LastId: lastID}, nil
	}
	return nil, nil
}

func validatePagination(resourceName, resourceID string, page *apiv1.Pagination) error {
	if page == nil {
		return nil
	}
	if page.Limit > 100 {
		return resource.NewValidationError(resourceName, resourceID, "max page limit is 100", nil)
	}
	return nil
}

func validateIncludes(resourceName string, includes []string) error {
	switch resourceName {
	case "job":
		for _, inc := range includes {
			if !backend.JobIncludes[inc] {
				return resource.NewValidationError(resourceName, "", fmt.Sprintf("invalid include %q", inc), nil)
			}
		}
	}
	return nil
}
