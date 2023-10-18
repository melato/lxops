package util

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type NamedNumber struct {
	Name  string
	Value int
}

func ReadNumbers(file string) ([]*NamedNumber, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(bytes.NewReader(content))
	var result []*NamedNumber
	for {
		fields, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(fields) != 2 {
			return nil, fmt.Errorf("expected 2 fields: %v", fields)
		}
		var num NamedNumber
		num.Name = fields[0]
		num.Value, err = strconv.Atoi(fields[1])
		if err != nil {
			return nil, err
		}
		result = append(result, &num)
	}
	return result, nil
}

func WriteNumbers(numbers []*NamedNumber, file string) error {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	for _, num := range numbers {
		writer.Write([]string{num.Name, strconv.Itoa(num.Value)})
	}
	writer.Flush()
	return os.WriteFile(file, buf.Bytes(), os.FileMode(0664))
}
