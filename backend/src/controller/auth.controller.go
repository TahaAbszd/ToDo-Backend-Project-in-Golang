package controller

import (
	"os"
	"time"
	"todo/src/db"
	"todo/src/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func RegisterUser(c *fiber.Ctx) error {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body request
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse"})
	}

	var existingUser models.User
	err := db.DB.Collection("users").FindOne(c.Context(), bson.M{
		"email": body.Email}).Decode(&existingUser)
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "useralerady exists",
		})
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "could not hash pass",
		})
	}
	newUser := models.User{
		Email:    body.Email,
		Password: string(hashedPassword),
	}
	_, err = db.DB.Collection("users").InsertOne(c.Context(), newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create user"})

	}
	return c.JSON(fiber.Map{"message": "User registered succesfully"})
}

func LoginUser(c *fiber.Ctx) error {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var body request
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse"})
	}
	var user models.User

	err := db.DB.Collection("users").FindOne(c.Context(), bson.M{"email": body.Email}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Credential"})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credential"})
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.Id.Hex(),
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	},
	)
	t, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not login"})
	}
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    t,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
	})
	return c.JSON(fiber.Map{"message": "LOgged in seccesfully"})
}

func LogoutUser(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-24 * time.Hour),
		HTTPOnly: true,
	})
	return c.JSON(fiber.Map{"message": "logged out"})
}
