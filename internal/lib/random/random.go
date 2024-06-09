package random

import (
	"math/rand"
	"time"
)

func NewRandomSTR(size int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	chars := []rune("SKLAJHFKJASDHFJKHASDJKF" + "kjashdfjkafnasdfas" + "0123456789")
	b := make([]rune, size)
	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))]
	}
	return string(b)
}
