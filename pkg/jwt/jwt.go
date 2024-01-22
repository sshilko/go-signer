package jwt

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type SimpleJWT struct {
	secret []byte
}

func NewSimpleJWT(secret []byte) *SimpleJWT {
	return &SimpleJWT{
		secret: secret,
	}
}

const userIDClaim = "sub"

func (d *SimpleJWT) ExampleToken(userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{userIDClaim: userID.String()})
	return token.SignedString(d.secret)
}

func (d *SimpleJWT) EchoMiddleware() (echo.MiddlewareFunc, error) {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			authorization := c.Request().Header.Get("Authorization")
			if authorization == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "No Authorization Header")
			}

			if !strings.HasPrefix(authorization, "Bearer ") {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Authorization Header")
			}

			token := strings.TrimPrefix(authorization, "Bearer ")

			claims := jwt.MapClaims{}
			_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
				return d.secret, nil
			})

			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Token")
			}

			c.Set("userID", claims[userIDClaim])
			return next(c)
		}
	}, nil
}
