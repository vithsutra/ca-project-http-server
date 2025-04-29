package utils

import (
	"errors"
	"os"
	"time"

	jwt_token "github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userId string, email string, userName string, adminId string) (string, error) {

	secretKey := os.Getenv("JWT_TOKEN_SCRETE_KEY")

	if secretKey == "" {
		return "", errors.New("missing JWT_TOKEN_SCRETE_KEY env variable")
	}

	token := jwt_token.NewWithClaims(
		jwt_token.SigningMethodHS256,
		jwt_token.MapClaims{
			"admin_id":  adminId,
			"id":        userId,
			"user_name": userName,
			"email":     email,
			"expiry":    time.Now().Add(time.Minute * 2).Unix(),
		},
	)

	tokenString, err := token.SignedString([]byte(secretKey))

	return tokenString, err
}
