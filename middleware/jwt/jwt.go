package jwt

import (
	"errors"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().URL.Path == "/health" {
				return next(c)
			}

			authHeader := c.Request().Header.Get("Authorization")
			// Browsers cannot set Authorization on a native WebSocket handshake.
			// Accept a token query parameter only for the notification WebSocket.
			if authHeader == "" && strings.HasSuffix(c.Request().URL.Path, "/notifications/ws") {
				if token := c.QueryParam("token"); token != "" {
					authHeader = "Bearer " + token
				}
			}
			if authHeader == "" {
				return echo.NewHTTPError(401, "Authorization header required")
			}

			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				return echo.NewHTTPError(401, "Invalid authorization header")
			}

			secretKey := os.Getenv("JWT_SECRET")
			if secretKey == "" {
				return echo.NewHTTPError(500, "JWT_SECRET is not configured")
			}

			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenParts[1], claims, func(token *jwt.Token) (interface{}, error) {
				if token.Method != jwt.SigningMethodHS256 {
					return nil, errors.New("unexpected JWT signing method")
				}
				return []byte(secretKey), nil
			}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

			if err != nil || token == nil || !token.Valid || claims.UserID <= 0 {
				return echo.NewHTTPError(401, "Invalid token")
			}

			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)
			c.Set("role", claims.Role)

			return next(c)
		}
	}
}
