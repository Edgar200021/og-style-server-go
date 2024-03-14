package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func SignJWT(payload jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	if tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET"))); err != nil {
		return "", err
	} else {
		return tokenString, err
	}

}

func ParseJWT(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims), token.Valid; ok {
		expiresDate, _ := time.Parse(time.RFC3339Nano, claims["expires"].(string))

		if time.Now().After(expiresDate) {
			return nil, errors.New("token expired")
		}
		return claims, nil
	} else {
		return nil, errors.New("something went wrong")
	}
}
