package utils

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

func SignUserJWT(secret, userID string) (string, error) {
	claims := jwt.MapClaims{
		"uid": userID,
		"exp": time.Now().Add(30 * 24 * time.Hour).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}

func ParseUserJWT(secret, token string) (string, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil || !t.Valid { return "", err }
	claims := t.Claims.(jwt.MapClaims)
	return claims["uid"].(string), nil
}
