package schema

import (
	"encoding/json"

	"github.com/jeffrom/job-manager/pkg/job"
	"github.com/qri-io/jsonschema"
)

func Parse(q *job.Queue) (*Schema, error) {
	var argSchema *jsonschema.Schema
	var dataSchema *jsonschema.Schema
	var resultSchema *jsonschema.Schema

	if argData := q.ArgSchemaRaw; len(argData) > 0 {
		argSchema = &jsonschema.Schema{}
		if err := json.Unmarshal(argData, argSchema); err != nil {
			return nil, err
		}
	}
	if dataData := q.DataSchemaRaw; len(dataData) > 0 {
		dataSchema = &jsonschema.Schema{}
		if err := json.Unmarshal(dataData, dataSchema); err != nil {
			return nil, err
		}
	}
	if resultData := q.ResultSchemaRaw; len(resultData) > 0 {
		resultSchema = &jsonschema.Schema{}
		if err := json.Unmarshal(resultData, resultSchema); err != nil {
			return nil, err
		}
	}

	return &Schema{
		Args:   argSchema,
		Data:   dataSchema,
		Result: resultSchema,
	}, nil
}

func ParseBytes(b []byte) (*Schema, error) {
	q := &job.Queue{}
	if err := json.Unmarshal(b, q); err != nil {
		return nil, err
	}
	return Parse(q)
}
