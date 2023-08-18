package middleware

import (
	"GoYin/server/service/api/models"
	"context"
	"errors"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/golang-jwt/jwt"
	"net/http"
)

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("couldn't handle this token")
	TokenNotFound    = errors.New("no token")
)

func JWTAuth(secretKey string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		token := c.Query("token")
		if token == "" {
			token = string(c.FormValue("token"))
			if token == "" {
				c.JSON(http.StatusInternalServerError, utils.H{
					"status_code": 500,
					"status_msg":  TokenNotFound.Error(),
				})
				c.Abort()
				return
			}
		}
		j := NewJWT(secretKey)
		// Parse the information contained in the token
		claims, err := j.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.H{
				"status_code": 500,
				"status_msg":  err.Error(),
			})
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Set("userId", claims.ID)
		c.Next(ctx)
	}
}

type JWT struct {
	SigningKey []byte
}

func NewJWT(secretKey string) *JWT {
	return &JWT{
		SigningKey: []byte(secretKey),
	}
}

// CreateToken to create a token
func (j *JWT) CreateToken(claims models.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// ParseToken to parse a token
func (j *JWT) ParseToken(tokenString string) (*models.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid
	} else {
		return nil, TokenInvalid
	}
}
