package auth

import (
	"crypto/rand"
	"strconv"
)

func randomDigit() int {
	randomByte := make([]byte, 1)
	_, err := rand.Read(randomByte)
	if err != nil {
		panic(err)
	}
	i := int(randomByte[0]) % 16
	if i < 10 {
		return i
	}
	return randomDigit()
}

func newLoginCode() string {
	s := ""
	for i := 0; i < 6; i++ {
		s += strconv.Itoa(randomDigit())
	}
	return s
}
