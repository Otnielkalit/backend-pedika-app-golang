package middleware

import (
	"backend-pedika-fiber/helper"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func AdminMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		response := helper.ResponseWithOutData{
			Code:    fiber.StatusUnauthorized,
			Status:  "error",
			Message: "Unauthorized: Missing token",
		}
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}
	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		response := helper.ResponseWithOutData{
			Code:    fiber.StatusUnauthorized,
			Status:  "error",
			Message: "Unauthorized: Invalid token format",
		}
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	tokenString := splitToken[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil || !token.Valid {
		response := helper.ResponseWithOutData{
			Code:    fiber.StatusUnauthorized,
			Status:  "error",
			Message: "Unauthorized: Invalid token",
		}
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	role := claims["role"].(string)
	if role != "admin" {
		response := helper.ResponseWithOutData{
			Code:    fiber.StatusForbidden,
			Status:  "error",
			Message: "Forbidden: Access Not Allowed",
		}
		return c.Status(fiber.StatusForbidden).JSON(response)
	}
	return c.Next()
}
