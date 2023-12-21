package util

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"time"
)

func JwtCreateToken(userID int64, jwtKey string, ttlDuration string) (tokenString string, err error) {
	now := time.Now()
	secretKey := []byte(jwtKey)
	ttl, _ := time.ParseDuration(ttlDuration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": userID,
			"iat": now.Unix(),
			"exp": now.Add(ttl).Unix(),
		},
	)

	tokenString, err = token.SignedString(secretKey)
	if err != nil {
		log.Println("error ", err)
		return
	}

	return
}

func JwtVerify(tokenString, jwtKey string) (isValid bool, err error) {
	secretKey := []byte(jwtKey)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return false, err
	}

	if !token.Valid {
		return false, fmt.Errorf("invalid token")
	}

	return true, nil
}

func JwtParseToken(tokenString, jwtKey string) (token *jwt.Token, err error) {
	secretKey := []byte(jwtKey)
	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return
	}

	return
}
