package schema

import (
	"encoding/json"
)

func Parse(b []byte) (*Schema, error) {
	if len(b) == 0 {
		return nil, nil
	}
	scm := &Schema{}
	if err := json.Unmarshal(b, scm); err != nil {
		return nil, err
	}
	return scm, nil
}
