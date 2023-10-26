package migrate

import (
	"strings"
)

func EqualYaml(d1, d2 []byte) bool {
	s1 := strings.TrimSpace(string(d1))
	s2 := strings.TrimSpace(string(d2))

	return s1 == s2
}
