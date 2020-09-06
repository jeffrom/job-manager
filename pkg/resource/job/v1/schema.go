package v1

import (
	"encoding/json"

	"github.com/jeffrom/job-manager/pkg/schema"
)

func ParseSchema(q *Queue) (*schema.Schema, error) {
	var scm *schema.Schema

	if b := q.Schema; len(b) > 0 {
		scm = &schema.Schema{}
		if err := json.Unmarshal(b, scm); err != nil {
			return nil, err
		}
	}

	return scm, nil
}
