package util

import (
	"github.com/golang-jwt/jwt/v5"
	"log"
	"swpr/config"
	"time"
)

func JwtCreateToken(userID int64) (tokenString string, err error) {
	now := time.Now()
	secretKey := []byte(config.Instance.Security.JwtSecKey)
	ttl, _ := time.ParseDuration(config.Instance.Security.JwtTTL)
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
