package util

import (
	"sort"
)

type StringSlice []string

func (t StringSlice) Equals(b []string) bool {
	if len(t) != len(b) {
		return false
	}
	for i, s := range t {
		if s != b[i] {
			return false
		}
	}
	return true
}

/** Return the elements of this slice that are not in the b slice */
func (t StringSlice) Diff(b []string) []string {
	var result []string
	set := StringSlice(b).ToSet()
	for _, s := range t {
		if !set.Contains(s) {
			result = append(result, s)
		}
	}
	return result
}

func (t StringSlice) Sorted() []string {
	result := make([]string, len(t))
	for i, s := range t {
		result[i] = s
	}
	sort.Strings(result)
	return result
}

func (t StringSlice) RemoveDuplicates() []string {
	return t.Union()
}

func (t StringSlice) Remove(remove string) []string {
	var removed bool
	var result []string
	for i, s := range t {
		if s == remove {
			if !removed {
				removed = true
				result = append(result, t[0:i]...)
			}
		} else if removed {
			result = append(result, s)

		}
	}
	if removed {
		return result
	} else {
		return t
	}
}

func (t StringSlice) ToSet() Set[string] {
	result := make(Set[string])
	for _, s := range t {
		result.Put(s)
	}
	return result
}

func (t StringSlice) Union(lists ...[]string) []string {
	set := make(Set[string])
	var result []string
	add := func(list []string) {
		for _, s := range list {
			if !set.Contains(s) {
				set.Put(s)
				result = append(result, s)
			}
		}
	}
	add(t)
	for _, list := range lists {
		add(list)
	}
	return result
}
