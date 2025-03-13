package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateOTP() (string, error) {
	otp := ""
	for i := 0; i < 4; i++ {
		// Generate a random number between 0 and 9
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		otp += fmt.Sprintf("%d", n.Int64())
	}
	return otp, nil
}
