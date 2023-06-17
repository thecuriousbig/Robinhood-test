package auth

import (
	"robinhood/config"

	"github.com/golang-jwt/jwt/v5"
)

type JWTCustomClaims struct {
	UserId string `json:"userId"`
	jwt.RegisteredClaims
}

func GenerateToken(UserId string) (string, error) {
	claims := JWTCustomClaims{
		UserId,
		jwt.RegisteredClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Get().JWT.Secret))
}

func ParseToken(tokenString string) (*JWTCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Get().JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JWTCustomClaims)
	if !ok {
		return nil, err
	}
	return claims, nil
}
