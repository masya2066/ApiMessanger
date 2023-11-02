package utils

import (
	"ApiMessenger/models"
	"github.com/dgrijalva/jwt-go"
)

func ParseToken(tokenString string) (claims *models.Claims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("$2a$14$cTtWCEgT1Vl2Q1orRe1GRObekuWfmW3DEhsf.I9kgoG45SkMg/Y.2"), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.Claims)

	if !ok {
		return nil, err
	}

	return claims, nil
}
