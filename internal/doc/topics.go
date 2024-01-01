package doc

import (
	"fmt"
	"io/fs"
	"path"
	"strings"
)

type Topics struct {
	FS     fs.FS  `name:"-"`
	Dir    string `name:"-"`
	Suffix string `name:"-"`
}

func NewTopics(fs fs.FS, dir string, suffix string) *Topics {
	return &Topics{FS: fs, Dir: dir, Suffix: suffix}
}

func (t *Topics) List() error {
	entries, err := fs.ReadDir(t.FS, t.Dir)
	if err != nil {
		return fmt.Errorf("ReadDir %s: %w", t.Dir, err)
	}
	for _, e := range entries {
		if !e.IsDir() {
			fname := e.Name()
			name := strings.TrimSuffix(fname, t.Suffix)
			if len(name) != len(fname) {
				fmt.Println(name)
			}
		}
	}
	return nil
}

func (t *Topics) PrintTopic(topic string) error {
	data, err := fs.ReadFile(t.FS, path.Join(t.Dir, topic+t.Suffix))
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", data)
	return nil
}

func (t *Topics) Print(topics ...string) error {
	switch len(topics) {
	case 0:
		return t.List()
	case 1:
		return t.PrintTopic(topics[0])
	default:
		return fmt.Errorf("select one topic")
	}
}
