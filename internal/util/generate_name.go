package util

import (
	"math/rand"
)

const alphabet string = "abcdefghijklmnopqrstuvwxyz"

func GenerateName(size int) string {
	buf := make([]byte, size)
	n := len(alphabet)
	for i := 0; i < size; i++ {
		k := rand.Intn(n)
		buf[i] = alphabet[k]
	}
	return string(buf)
}
