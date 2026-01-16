package templatefs

import (
	"bytes"
	"io"
	"io/fs"
	"strings"
	"text/template"
)

// TemplateFS - An FS that wraps another FS of templates.
// When it reads a file, it executes it as a template
// and returns the result.
type TemplateFS struct {
	FS           fs.FS
	TemplateData any
}

func NewTemplateFS(fsys fs.FS, templateData any) *TemplateFS {
	return &TemplateFS{FS: fsys, TemplateData: templateData}
}

func ExecuteTemplate(s string, model any) ([]byte, error) {
	tpl, err := template.New("").Parse(s)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf, model)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *TemplateFS) ReadFrom(r io.Reader) ([]byte, error) {
	var builder strings.Builder
	_, err := io.Copy(&builder, r)
	if err != nil {
		return nil, err
	}
	return ExecuteTemplate(builder.String(), t.TemplateData)
}

func (t *TemplateFS) ReadFile(name string) ([]byte, error) {
	f, err := t.FS.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return t.ReadFrom(f)
}

func (t *TemplateFS) Open(name string) (fs.File, error) {
	return t.FS.Open(name)
}
