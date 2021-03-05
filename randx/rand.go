package randx

import (
	"bytes"
	"crypto/rand"
	"math/big"
	rand2 "math/rand"
)

func RandomString(len int) string {
	var container string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(str[randomInt.Int64()])
	}
	return container
}

func Int(n int) int {
	if n <= 0 {
		return n
	}
	return rand2.Intn(n)
}
