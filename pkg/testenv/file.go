package testenv

import (
	"os"
	"testing"
)

func ReadFile(t testing.TB, p string) []byte {
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
