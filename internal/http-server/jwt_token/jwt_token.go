package jwt_token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
	"time"
)

const (
	tokenTTL = 12 * time.Hour

	// TODO: вынести ключ в переменнцю окружения
	signedKey = "very_secret_key"
)

type claimsWithId struct {
	jwt.StandardClaims
	Id int `json:"id"`
}

func GenerateToken(id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claimsWithId{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Id: id,
	})

	return token.SignedString([]byte(signedKey))
}

func ParseToken(t string) (int, error) {
	token, err := jwt.ParseWithClaims(t, &claimsWithId{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signin method")
		}
		return []byte(signedKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*claimsWithId)
	if !ok {
		return 0, fmt.Errorf("invalid token claims")
	}

	return claims.Id, nil
}

func GetTokenFromRequest(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	authParts := strings.Split(auth, " ")
	if len(authParts) != 2 {
		return "", fmt.Errorf("invalid token")
	}
	return authParts[1], nil
}
