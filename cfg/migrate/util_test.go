package migrate

import (
	"embed"
	"testing"
)

//go:embed test/*.yaml
var testFS embed.FS

func TestEqualYaml(t *testing.T) {
	a := []byte("a")
	b := []byte("a ")
	if !EqualYaml(a, b) {
		t.Fail()
	}
}
