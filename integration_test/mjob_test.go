package integration

import (
	"testing"

	"github.com/jeffrom/job-manager/mjob/consumer"
)

// TODO remove this?
func TestIntegrationMJob(t *testing.T) {
	tcs := []struct {
		name string
		cfg  consumer.Config
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
