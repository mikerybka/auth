package auth

import (
	"strconv"
	"time"
)

func newUserID() string {
	n := time.Now().UnixMilli() - 1737227100000
	return "user" + strconv.FormatInt(n, 10)
}
