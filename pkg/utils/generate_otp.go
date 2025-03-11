package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateOTP() (string, error) {
	max := big.NewInt(1000000) // 6 digits => 000000 to 999999

	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%06d", n.Int64()), nil
}
