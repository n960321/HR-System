package jwthelper

import (
	"HRSystem/internal/model"
	"HRSystem/pkg/errors"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/golang-jwt/jwt"
)

var JWTHelper *helper
var once sync.Once

const (
	ACCOUNT_TOKEN_KEY = "account:%v:token"
)

type helper struct {
	redisClient *redis.Client
}

var secretKey = []byte("just-for-practice")

type Claim struct {
	ID          uint64
	Account     string
	AccountType model.AccountType
	jwt.StandardClaims
}

func New(redisClient *redis.Client) {
	once.Do(func() {
		JWTHelper = &helper{
			redisClient: redisClient,
		}
	})

}

func (j *helper) GenerateJWTToken(id uint64, account string, accountType model.AccountType) (string, error) {
	expiresAt := time.Now().Add(1 * time.Hour).Unix()
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
	if j.redisClient != nil {
		err = j.redisClient.Set(fmt.Sprintf(ACCOUNT_TOKEN_KEY, id), tokenString, 1*time.Hour).Err()
		if err != nil {
			return "", err
		}
	}

	return tokenString, nil
}

func (j *helper) ParseToken(tokenString string) (*Claim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claim{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err == nil && token != nil {
		if claim, ok := token.Claims.(*Claim); ok && token.Valid {
			redisToken, err := j.redisClient.Get(fmt.Sprintf(ACCOUNT_TOKEN_KEY, claim.ID)).Result()
			if err != nil {
				return nil, err
			}

			if tokenString != redisToken {
				return nil, errors.ErrInvalidToken
			}
			return claim, nil
		}
	}
	return nil, err
}

func (j *helper) GetClaim(ctx *gin.Context) (*Claim, error) {
	token := ctx.Request.Header.Get("token")
	if token == "" {
		return nil, fmt.Errorf("no jwt token")
	}

	claim, err := j.ParseToken(token)
	if err != nil {
		return nil, fmt.Errorf("bad jwt: %s", err)
	}

	return claim, nil
}
