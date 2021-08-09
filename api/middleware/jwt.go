package middleware

import (
	"cloud-batch/configs"
	"cloud-batch/internal/pkg/app"
	"cloud-batch/internal/pkg/e"
	"cloud-batch/internal/pkg/logging"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"time"
)

var jwtSecret = []byte(configs.Server.JwtSecret)

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

func GenerateToken(username, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)

	claims := Claims{
		username,
		password,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "cloud-batch",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

func ParseToken(token string) (*Claims, error) {
	logging.Logger.Debugf("parse token: %s", token)
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := app.Gin{C: c}
		code := e.StatusOK
		token := c.GetHeader("Authorization")
		if token == "" {
			code = e.InvalidToken
		} else {
			claims, err := ParseToken(token)
			if err != nil {
				code = e.InvalidToken
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = e.InvalidToken
			}
		}

		if code != e.StatusOK {
			appG.ResponseI18nMsg(code, nil, nil, nil, nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
