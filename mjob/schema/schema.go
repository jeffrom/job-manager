// Package schema validates job arguments and result data with jsonschema.
package schema

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/hashicorp/go-multierror"
	"github.com/qri-io/jsonschema"
	jsonmin "github.com/tdewolff/minify/v2/json"
)

// SelfSchema is used to validate job schemas themselves.
var SelfSchema = &jsonschema.Schema{}

func init() {
	min := jsonmin.DefaultMinifier
	min.Precision = 0

	if err := json.Unmarshal(selfSchemaRaw, SelfSchema); err != nil {
		panic(err)
	}
}

// Schema provides jsonschema validation for job arguments, data, and result
// data.
type Schema struct {
	Args   *jsonschema.Schema `json:"args,omitempty"`
	Data   *jsonschema.Schema `json:"data,omitempty"`
	Result *jsonschema.Schema `json:"result,omitempty"`
}

func (s *Schema) Validate(ctx context.Context, args, data, result []byte) error {
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

func (s *Schema) ValidateArgs(ctx context.Context, jsonData []byte) error {
	if s == nil || s.Args == nil {
		return nil
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

func (s *Schema) ValidateResult(ctx context.Context, jsonData []byte) error {
	if s == nil || s.Result == nil {
		return nil
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

func (s *Schema) ValidateData(ctx context.Context, jsonData []byte) error {
	if s == nil || s.Data == nil {
		return nil
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

// Canonicalize deterministically formats json so it can be reliably compared
// to previous versions.
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
