package main

import (
	"errors"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func ExtractEmail(token string) (error, interface{}) {
	claims := jwt.MapClaims{}
	extractedToken := strings.Split(token, "Bearer ")
	if len(extractedToken) == 2 {
		token = strings.TrimSpace(extractedToken[1])
	} else {
		return errors.New("Incorrect Format of Authorization Token"), ""
	}
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SIGNING_KEY")), nil
	})

	if err != nil {
		return errors.New(err.Error()), ""
	}

	return nil, claims["Email"]
}