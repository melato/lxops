package util

import (
	"bufio"
	"io"
	"os"
	"strings"
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

func ReadProperties(file string) (map[string]string, error) {
	var f, err = os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	props := make(map[string]string)
	var reader = bufio.NewReader(f)
	for {
		var line, err = reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		sep := strings.IndexAny(line, "=:")
		if sep >= 0 {
			key := strings.TrimSpace(line[0:sep])
			value := strings.TrimSpace(line[sep+1:])
			props[key] = value
		}
	}
	return props, nil

}

func EscapeShell(args ...string) string {
	var buf strings.Builder
	for i, arg := range args {
		if i > 0 {
			buf.WriteString(" ")
		}
		if arg == "" {
			buf.WriteString("''")
		} else if strings.Contains(arg, " ") {
			buf.WriteString("\"")
			buf.WriteString(arg)
			buf.WriteString("\"")
		} else {
			buf.WriteString(arg)
		}
	}
	return buf.String()
}
