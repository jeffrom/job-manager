package job

import (
	"encoding/json"

	"github.com/jeffrom/job-manager/pkg/schema"
	"github.com/qri-io/jsonschema"
)

func Parse(q *Queue) (*schema.Schema, error) {
	var argSchema *jsonschema.Schema
	var dataSchema *jsonschema.Schema
	var resultSchema *jsonschema.Schema

	if argData := q.Schema.Args; len(argData) > 0 {
		argSchema = &jsonschema.Schema{}
		if err := json.Unmarshal(argData, argSchema); err != nil {
			return nil, err
		}
	}
	if dataData := q.Schema.Data; len(dataData) > 0 {
		dataSchema = &jsonschema.Schema{}
		if err := json.Unmarshal(dataData, dataSchema); err != nil {
			return nil, err
		}
	}
	if resultData := q.Schema.Result; len(resultData) > 0 {
		resultSchema = &jsonschema.Schema{}
		if err := json.Unmarshal(resultData, resultSchema); err != nil {
			return nil, err
		}
	}

	return &schema.Schema{
		Args:   argSchema,
		Data:   dataSchema,
		Result: resultSchema,
	}, nil
}

func ParseBytes(b []byte) (*schema.Schema, error) {
	q := &Queue{}
	if err := json.Unmarshal(b, q); err != nil {
		return nil, err
	}
	return Parse(q)
}
