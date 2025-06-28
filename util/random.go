package util

import (
	"math/rand"
	"strings"
	"time"
)

var seed *rand.Rand

func init() {
	//rand.Seed(time.Now().UnixNano())
	seed = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomInteger(max int64) int64 {
	return seed.Int63n(max)
}

func RandomString(n int) string {
	chars := "abcdefghijklmnopqrstuvwxyz"
	var str strings.Builder
	for i := 0; i < n; i++ {
		idx := int(RandomInteger(int64(len(chars))))
		str.WriteByte(chars[idx])
	}
	return str.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomAmount() int64 {
	return RandomInteger(1000)
}

func RandomCurrency() string {
	curr := []string{"USD", "Euros", "Rupees"}
	n := len(curr)
	return curr[seed.Intn(n)]
}
