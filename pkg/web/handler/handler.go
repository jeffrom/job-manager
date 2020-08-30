// Package handler contains http handlers.
package handler

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"google.golang.org/protobuf/proto"

	"github.com/jeffrom/job-manager/pkg/backend"
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
			log.Printf("error %T: %+v", err, err)
			if herr, ok := err.(httpError); ok {
				w.WriteHeader(herr.Status())
			} else if errors.Is(err, io.ErrUnexpectedEOF) {
				w.WriteHeader(http.StatusBadRequest)
			} else if errors.Is(err, backend.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}

			if pr, ok := err.(protoError); ok {
				if err := MarshalResponse(w, r, pr.Message()); err != nil {
					log.Printf("handler: error marshalling response: %v", err)
				}
			}
		}
	}
}

func UnmarshalBody(r *http.Request, v interface{}, required bool) error {
	defer r.Body.Close()
	ct := r.Header.Get("content-type")
	log.Printf("content type: %q", ct)

	b, rerr := ioutil.ReadAll(r.Body)
	if rerr != nil {
		return rerr
	}
	if len(b) == 0 {
		if required {
			return io.ErrUnexpectedEOF
		}
		return nil
	}

	var err error
	switch ct {
	case "", "application/json":
		err = json.Unmarshal(b, v)
	case "application/protobuf":
		err = proto.Unmarshal(b, v.(proto.Message))
	default:
		panic("handler: unknown content-type: " + ct)
	}

	log.Printf("params (%T): %+v", v, v)
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

	log.Printf("response (%T): %+v", v, v)
	if _, err := w.Write(b); err != nil {
		return err
	}
	return err
}
