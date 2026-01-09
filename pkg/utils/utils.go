package utils

import (
	crand "crypto/rand"
	"math/big"
	"time"
)

func GetTimestamp() int64 {
	return time.Now().Unix()
}

func GenerateKey() (string, error) {
	//rand.Seed(time.Now().UnixNano())
	return GenerateRandomCharsKey(48)
}

const keyChars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateRandomCharsKey(length int) (string, error) {
	b := make([]byte, length)
	maxI := big.NewInt(int64(len(keyChars)))

	for i := range b {
		n, err := crand.Int(crand.Reader, maxI)
		if err != nil {
			return "", err
		}
		b[i] = keyChars[n.Int64()]
	}

	return string(b), nil
}
