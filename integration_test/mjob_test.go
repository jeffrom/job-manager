package integration

import (
	"testing"

	"github.com/jeffrom/job-manager/mjob"
)

func TestIntegrationMJob(t *testing.T) {
	tcs := []struct {
		name string
		cfg  mjob.ConsumerConfig
	}{
		{
			name: "default",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

		})
	}
}
