package jwthelper

import (
	"HRSystem/internal/model"
	"time"

	"github.com/golang-jwt/jwt"
)

var secretKey = []byte("just-for-practice")

type Claim struct {
	ID          uint64
	Account     string
	AccountType model.AccountType
	jwt.StandardClaims
}

func GenerateJWTToken(id uint64, account string, accountType model.AccountType) (string, error) {
	expiresAt := time.Now().Add(10 * time.Second).Unix()
	claims := Claim{
		ID:          id,
		Account:     account,
		AccountType: accountType,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "admin",
			ExpiresAt: expiresAt,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*Claim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claim{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err == nil && token != nil {
		if claim, ok := token.Claims.(*Claim); ok && token.Valid {
			return claim, nil
		}
	}
	return nil, err
}
