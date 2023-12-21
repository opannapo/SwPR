package util

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"strconv"
	"time"
)

func JwtCreateToken(userID int64, jwtKey string, ttlDuration string) (tokenString string, err error) {
	now := time.Now()
	secretKey := []byte(jwtKey)
	ttl, _ := time.ParseDuration(ttlDuration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": strconv.Itoa(int(userID)),
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
		log.Println("error JwtParseToken", err)
		return
	}

	return
}

func JwtGetSubjectUserID(jwt, jwtKey string) (userID int, err error) {
	token, err := JwtParseToken(jwt, jwtKey)
	if err != nil {
		log.Println("error JwtParseToken", err)
		return
	}

	userIdStr, err := token.Claims.GetSubject()
	if err != nil {
		log.Println("error Claims.GetSubject()", err)
		return
	}
	userID, _ = strconv.Atoi(userIdStr)
	return
}
