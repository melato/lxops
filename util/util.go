package util

import (
	"os"
)

func FileExists(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		return false
	}
	return true
}

func DirExists(dir string) bool {
	st, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return st.IsDir()
}
