package lxops

import (
	"sort"

	"strings"
)

type Path string

// IsDescendantOf returns true if c is under p's filesystem hierarchy
func (c Path) IsDescendantOf(p string) bool {
	if p == "" {
		return false
	}
	if !strings.HasPrefix(string(c), p) {
		return false
	}
	if len(p) == len(c) {
		return true
	}
	return c[len(p)] == '/'
}

// RootPaths returns a copy of paths , removing any path that is a descendant of another one
func RootPaths(paths []string) []string {
	sorted := make([]string, len(paths))
	copy(sorted, paths)
	sort.Strings(sorted)
	var roots []string
	var last string
	for _, s := range sorted {
		if !Path(s).IsDescendantOf(last) {
			roots = append(roots, s)
			last = s
		}
	}
	return roots
}
