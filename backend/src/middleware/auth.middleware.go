package middleware

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func AuthMiddleware(c *fiber.Ctx) error {
	tokenString := c.Cookies("jwt")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unathorized",
		})
	}
	claims := token.Claims.(jwt.MapClaims)
	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "token expired",
		})
	}
	c.Locals("userId", claims["userId"])
	return c.Next()
}
