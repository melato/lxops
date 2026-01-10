package cfg

import (
	"melato.org/lxops/util"
)

type GetVariable func(name string) (string, bool)

/*
Modify and filter a path.

Return the effective path and a boolean indicating whether it should be used or not.

path goes through variable substitution.
If the substitution fails, the path is ignored (returns false).
This provides a simple way of conditional configuration.
Just use a variable in the path that can be defined or not.
*/
func filterPath(path HostPath, getVariable GetVariable) (HostPath, bool) {
	file := string(path)
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
