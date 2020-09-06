package schema

import (
	"encoding/json"
)

func Parse(b []byte) (*Schema, error) {
	scm := &Schema{}
	if err := json.Unmarshal(b, scm); err != nil {
		return nil, err
	}
	return scm, nil
}
