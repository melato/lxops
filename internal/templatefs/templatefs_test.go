package templatefs

import (
	"embed"
	"io/fs"
	"testing"
)

//go:embed test/a.tpl
var testFS embed.FS

func TestTemplateFS(t *testing.T) {
	data := map[string]string{
		"A": "a",
	}
	fsys := NewTemplateFS(testFS, data)
	b, err := fs.ReadFile(fsys, "test/a.tpl")
	if err != nil {
		t.Fatalf("%v", err)
	}
	s := string(b)
	if s != "a" {
		t.Fatalf("%s", s)
	}
}
