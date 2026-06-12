package middleware

import (
	"strings"

	"github.com/Victor-132/cashtrackr/internal/auth"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")

	if authHeader == "" {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	userID, err := auth.ValidateJWT(tokenString)
	if err != nil {
		return err
	}

	ctx.Locals("user_id", userID)

	return ctx.Next()
}
