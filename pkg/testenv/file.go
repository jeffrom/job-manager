package testenv

import (
	"io/ioutil"
	"testing"
)

func ReadFile(t testing.TB, p string) []byte {
	b, err := ioutil.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
