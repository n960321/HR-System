package middleware

import (
	"HRSystem/pkg/errors"
	"HRSystem/pkg/jwthelper"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func ValidateJWT() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("token")
		if token == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": errors.ErrInvalidToken.Error()})
			ctx.Abort()
			return
		}

		_, err := jwthelper.ParseToken(token)
		if err != nil {
			innerErr := errors.ErrInvalidToken
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorExpired != 0 {
					innerErr = errors.ErrTokenExpired
				}
			}
			ctx.JSON(http.StatusBadRequest, gin.H{"message": innerErr.Error()})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
