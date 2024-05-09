package jwt_token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"time"
)

const (
	tokenTTL = 12 * time.Hour

	// TODO: вынести ключ в переменнцю окружения
	signedKey = "very_secret_key"
)

const (
	RoleAdmin    = "admin"
	RoleCustomer = "customer"
)

type extendedClaims struct {
	jwt.StandardClaims
	Id    int    `json:"id"`
	Role  string `json:"role"`
	Alias string `json:"alias"`
}

func GenerateToken(id int, role, alias string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &extendedClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Id:    id,
		Role:  role,
		Alias: alias,
	})

	return token.SignedString([]byte(signedKey))
}

func ParseToken(t string) (id int, role, alias string, err error) {
	token, err := jwt.ParseWithClaims(t, &extendedClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signin method")
		}
		return []byte(signedKey), nil
	})
	if err != nil {
		return 0, "", "", err
	}

	claims, ok := token.Claims.(*extendedClaims)
	if !ok {
		return 0, "", "", fmt.Errorf("invalid token claims")
	}

	return claims.Id, claims.Role, claims.Alias, nil
}

func GetTokenFromRequest(auth string) (string, error) {
	authParts := strings.Split(auth, " ")
	if len(authParts) != 2 {
		return "", fmt.Errorf("invalid token")
	}
	return authParts[1], nil
}
