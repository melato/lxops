package templatefs

import (
	"bytes"
	"io"
	"io/fs"
	"strings"
	"text/template"
	"time"
)

// TemplateFS - An FS that wraps another FS of templates.
// When it reads a file, it executes the template
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
	data, err := t.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return NewFile(name, data), nil
}

type bytesFile struct {
	Filename string
	Filesize int64
	Reader   io.Reader
}

func NewFile(name string, b []byte) fs.File {
	return &bytesFile{
		Filename: name,
		Filesize: int64(len(b)),
		Reader:   bytes.NewReader(b),
	}
}

func (t *bytesFile) Stat() (fs.FileInfo, error) {
	return t, nil
}
func (t *bytesFile) Read(b []byte) (int, error) {
	return t.Reader.Read(b)
}
func (t *bytesFile) Close() error {
	return nil
}

func (t *bytesFile) Name() string       { return t.Filename }
func (t *bytesFile) Size() int64        { return t.Filesize }
func (t *bytesFile) Mode() fs.FileMode  { return fs.FileMode(0555) }
func (t *bytesFile) ModTime() time.Time { return time.Now() }
func (t *bytesFile) IsDir() bool        { return false }
func (t *bytesFile) Sys() any           { return nil }
