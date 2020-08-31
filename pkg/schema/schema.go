// Package schema contains code for dealing with jsonschema.
package schema

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/qri-io/jsonschema"
	jsonmin "github.com/tdewolff/minify/v2/json"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var ArgSchema *jsonschema.Schema = &jsonschema.Schema{}
var DataSchema *jsonschema.Schema = &jsonschema.Schema{}
var ResultSchema *jsonschema.Schema = &jsonschema.Schema{}

func init() {
	min := jsonmin.DefaultMinifier
	min.Precision = 0

	if err := json.Unmarshal(argSchemaRaw, ArgSchema); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(dataSchemaRaw, DataSchema); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(resultSchemaRaw, ResultSchema); err != nil {
		panic(err)
	}
}

type Schema struct {
	Args   *jsonschema.Schema
	Data   *jsonschema.Schema
	Result *jsonschema.Schema
}

// func (s *Schema) Validate(ctx context.Context,

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
