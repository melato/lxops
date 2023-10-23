package lxops

import (
	"sort"
	"strings"

	"melato.org/lxops/cfg"
)

type InstanceFS struct {
	Id         string
	Path       string
	IsNew      bool
	Filesystem *cfg.Filesystem
}

func (t *InstanceFS) IsDir() bool {
	return strings.HasPrefix(t.Path, "/")
}

func (t *InstanceFS) IsZfs() bool {
	return !t.IsDir()
}

func (t *InstanceFS) Dir() string {
	if t.IsDir() {
		return string(t.Path)
	} else {
		return "/" + string(t.Path)
	}
}

type InstanceFSList []*InstanceFS

func (t InstanceFSList) Len() int           { return len(t) }
func (t InstanceFSList) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t InstanceFSList) Less(i, j int) bool { return t[i].Path < t[j].Path }

func (t InstanceFSList) Sort() { sort.Sort(t) }

// Roots returns the FSPaths that are not children of other FSPaths, within this set.
func (t InstanceFSList) Roots() []*InstanceFS {
	paths := make([]string, t.Len())
	fsMap := make(map[string]int)
	for i, fs := range t {
		paths[i] = fs.Path
		fsMap[fs.Path] = i
	}
	rootPaths := RootPaths(paths)
	roots := make([]*InstanceFS, len(rootPaths))
	for i, path := range rootPaths {
		roots[i] = t[fsMap[path]]
	}
	return roots
}
