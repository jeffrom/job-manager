// Package handler contains http handlers.
package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"google.golang.org/protobuf/proto"

	"github.com/jeffrom/job-manager/pkg/backend"
	"github.com/jeffrom/job-manager/pkg/schema"
	"github.com/jeffrom/job-manager/pkg/web/middleware"
)

type httpError interface {
	error
	Status() int
}

type protoError interface {
	Message() proto.Message
}

func Func(fn func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			reqLog := middleware.RequestLogFromContext(r.Context())

			// set status depending on the error returned
			status := http.StatusInternalServerError
			if herr, ok := err.(httpError); ok {
				status = herr.Status()
			} else if errors.Is(err, io.ErrUnexpectedEOF) {
				status = http.StatusBadRequest
			} else if errors.Is(err, backend.ErrNotFound) {
				status = http.StatusNotFound
			} else if errors.Is(err, &backend.VersionConflictError{}) {
				status = http.StatusConflict
			} else if errors.Is(err, &schema.ValidationError{}) {
				status = http.StatusBadRequest
			}

			reqLog.Str("err_type", fmt.Sprintf("%T", err))
			if status >= http.StatusInternalServerError {
				reqLog.Err(err)
			}

			// maybe translate errors into api errors
			vcErr := &backend.VersionConflictError{}
			if errors.As(err, &vcErr) {
				err = newVersionConflictError(vcErr)
			}

			// log the error
			log := middleware.LoggerFromContext(r.Context())
			logRequestError(log, status, err)
			w.WriteHeader(status)

			if pr, ok := err.(protoError); ok {
				if err := MarshalResponse(w, r, pr.Message()); err != nil {
					middleware.LoggerFromContext(r.Context()).Error().Err(err).Msg("marshal response failed")
				}
			}
		}
	}
}

func UnmarshalBody(r *http.Request, v interface{}, required bool) error {
	log := middleware.LoggerFromContext(r.Context())
	defer r.Body.Close()
	ct := r.Header.Get("content-type")
	var b []byte
	if r.Method == "GET" {
		ct = "application/x-www-form-urlencoded"
	} else {
		var rerr error
		b, rerr = ioutil.ReadAll(r.Body)
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

	log := middleware.LoggerFromContext(r.Context())
	log.Debug().
		Interface("data", v).
		Str("type", fmt.Sprintf("%T", v)).
		Msg("response")
	if _, err := w.Write(b); err != nil {
		return err
	}
	return err
}

func logRequestError(log *middleware.Logger, status int, err error) {
	ev := log.Warn()
	if status >= http.StatusInternalServerError {
		ev = log.Error()
	}

	ev.Err(err)

	ev.Msg("request error")
}
