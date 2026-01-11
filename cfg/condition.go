package cfg

import (
	"fmt"
	"strings"

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
func filterPath(path HostPath, getVariable GetVariable) (HostPath, bool, error) {
	file := string(path)
	reserved := strings.IndexAny(file, "|;,:")
	if reserved > 0 {
		return "", false, fmt.Errorf("path %s contains reserved character '%v'", file[reserved])
	}
	file, err := util.Substitute(file, getVariable)
	if err == nil {
		return HostPath(file), true, nil
	} else {
		return "", false, nil
	}
}

func filterPaths(paths []HostPath, getVariable GetVariable) ([]HostPath, error) {
	var k int
	for _, path := range paths {
		path, pass, err := filterPath(path, getVariable)
		if err != nil {
			return nil, err
		}
		if pass {
			paths[k] = path
			k++
		}
	}
	if k == len(paths) {
		return paths, nil
	} else {
		return paths[0:k], nil
	}
}
