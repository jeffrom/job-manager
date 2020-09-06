// Package schema contains code for dealing with jsonschema.
package schema

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/hashicorp/go-multierror"
	"github.com/qri-io/jsonschema"
	jsonmin "github.com/tdewolff/minify/v2/json"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var SelfSchema *jsonschema.Schema = &jsonschema.Schema{}

func init() {
	min := jsonmin.DefaultMinifier
	min.Precision = 0

	if err := json.Unmarshal(selfSchemaRaw, SelfSchema); err != nil {
		panic(err)
	}
}

type Schema struct {
	Args   *jsonschema.Schema `json:"args,omitempty"`
	Data   *jsonschema.Schema `json:"data,omitempty"`
	Result *jsonschema.Schema `json:"result,omitempty"`
}

func (s *Schema) Validate(ctx context.Context, args, data, result interface{}) error {
	merr := &multierror.Error{}
	if err := s.ValidateArgs(ctx, args); err != nil {
		merr = multierror.Append(merr, err)
	}
	if err := s.ValidateData(ctx, data); err != nil {
		merr = multierror.Append(merr, err)
	}
	if err := s.ValidateResult(ctx, result); err != nil {
		merr = multierror.Append(merr, err)
	}
	return merr.ErrorOrNil()
}

func (s *Schema) ValidateArgs(ctx context.Context, arg interface{}) error {
	if s.Args == nil {
		return nil
	}

	jsonData, err := marshalToJSON(arg)
	if err != nil {
		return err
	}
	keyErrs, err := s.Args.ValidateBytes(ctx, jsonData)
	if err != nil {
		return err
	}
	if len(keyErrs) > 0 {
		return NewValidationErrorKeyErrs(keyErrs)
	}
	return nil
}

func (s *Schema) ValidateResult(ctx context.Context, arg interface{}) error {
	if s.Result == nil {
		return nil
	}

	jsonData, err := marshalToJSON(arg)
	if err != nil {
		return err
	}

	keyErrs, err := s.Result.ValidateBytes(ctx, jsonData)
	if err != nil {
		return err
	}
	if len(keyErrs) > 0 {
		return NewValidationErrorKeyErrs(keyErrs)
	}
	return nil
}

func (s *Schema) ValidateData(ctx context.Context, arg interface{}) error {
	if s.Data == nil {
		return nil
	}

	jsonData, err := marshalToJSON(arg)
	if err != nil {
		return err
	}

	keyErrs, err := s.Result.ValidateBytes(ctx, jsonData)
	if err != nil {
		return err
	}
	if len(keyErrs) > 0 {
		return NewValidationErrorKeyErrs(keyErrs)
	}
	return nil
}

func marshalToJSON(arg interface{}) ([]byte, error) {
	var err error
	var jsonData []byte
	if msg, ok := arg.(proto.Message); ok {
		jsonData, err = protojson.Marshal(msg)
	} else {
		jsonData, err = json.Marshal(arg)
	}
	return jsonData, err
}

func Canonicalize(data []byte) ([]byte, error) {
	b := &bytes.Buffer{}
	if err := jsonmin.Minify(nil, b, bytes.NewReader(data), nil); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// func Equals(a, b *Schema) bool {
// 	if a == nil || b == nil {
// 		return a == b
// 	}
// 	// TODO quick solution is to serialize em to compare
// 	return false
// }
