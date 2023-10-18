package yaml

func FirstLineIs(data []byte, line string) bool {
	n := len(line)
	if len(data) < n {
		return false
	}
	if len(data) > n {
		c := rune(data[n])
		if (c != '\r') && (c != '\n') {
			return false
		}
	}
	for i := 0; i < n; i++ {
		if data[i] != line[i] {
			return false
		}
	}
	return true
}

func FirstLineComment(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	if data[0] != '#' {
		return ""
	}
	for i, c := range data {
		if (c == '\r') || (c == '\n') {
			return string(data[0:i])
		}
	}
	return string(data)
}
