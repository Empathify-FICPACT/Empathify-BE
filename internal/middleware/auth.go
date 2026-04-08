package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"

	"github.com/Empathify-FICPACT/Empathify-BE/pkg/jwt"
	"github.com/Empathify-FICPACT/Empathify-BE/pkg/response"
)

const UserIDKey = "user_id"

func AuthMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		var tokenStr string

		cookieToken := c.Cookies("access_token")
		if cookieToken != "" {
			tokenStr = cookieToken
		} else {
			authHeader := c.Get("Authorization")
			if authHeader == "" {
				return response.Unauthorized(c, "missing token")
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return response.Unauthorized(c, "invalid authorization format")
			}

			tokenStr = parts[1]
		}

		claims, err := jwt.ParseToken(tokenStr)
		if err != nil {
			switch err {
			case jwt.ErrExpiredToken:
				return response.Unauthorized(c, "token has expired")
			default:
				return response.Unauthorized(c, "invalid token")
			}
		}

		c.Locals(UserIDKey, claims.UserID)

		return c.Next()
	}
}

// helper buat ambil user_id di handler
func GetUserID(c fiber.Ctx) string {
	return c.Locals(UserIDKey).(string)
}