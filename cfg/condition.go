package cfg

import (
	"strings"

	"melato.org/lxops/util"
)

type GetVariable func(name string) (string, bool)

/*
	 paths have the form
		path
		variable|path
*/
func filterPath(path HostPath, getVariable GetVariable) (HostPath, bool) {
	name, file, hasCondition := strings.Cut(string(path), "|")
	if hasCondition {
		_, hasVariable := getVariable(name)
		if !hasVariable {
			return "", false
		}
	} else {
		file = string(path)
	}
	file, err := util.Substitute(file, getVariable)
	if err == nil {
		return HostPath(file), true
	} else {
		return "", false
	}
}

func filterPaths(paths []HostPath, getVariable GetVariable) []HostPath {
	var k int
	for _, path := range paths {
		path, pass := filterPath(path, getVariable)
		if pass {
			paths[k] = path
			k++
		}
	}
	if k == len(paths) {
		return paths
	} else {
		return paths[0:k]
	}
}
