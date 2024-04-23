package jwt_token

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	tokenTTL  = 12 * time.Hour
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
