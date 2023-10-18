package lxops

import (
	"math/rand"
	"strconv"
	"time"
)

var seeded bool

func RandomDeviceName() string {
	if !seeded {
		rand.Seed(time.Now().Unix())
		seeded = true
	}
	n := rand.Int31()
	return "d" + strconv.FormatInt(int64(n), 36)
}
