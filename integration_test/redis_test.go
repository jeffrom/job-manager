package integration

import (
	"fmt"
	"testing"

	"github.com/jeffrom/job-manager/pkg/backend/beredis"
)

func TestIntegrationRedis(t *testing.T) {
	be := beredis.New()
	fmt.Println(be)
}
