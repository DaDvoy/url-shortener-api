package random

import (
	"math/rand"
	"time"
)

func New() string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"0123456789_"

	randStr := make([]byte, 10)
	for i := range randStr {
		randStr[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(randStr)
}
