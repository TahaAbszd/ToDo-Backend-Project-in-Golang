package controller

import (
	"todo/src/db"
	"todo/src/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateTodo(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)
	type body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}
	var data body
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse json",
		})
	}
	if data.Status != string(models.StatusCompleted) &&
		data.Status != string(models.StatusINcomplete) {
		data.Status = string(models.StatusINcomplete)
	}
	todo := bson.M{
		"_id":         primitive.NewObjectID(),
		"title":       data.Title,
		"description": data.Description,
		"status":      data.Status,
		"userId":      userId,
	}

	_, err := db.DB.Collection("todos").InsertOne(c.Context(), todo)
	if err != nil {
		return c.Status(fiber.StatusInsufficientStorage).JSON(fiber.Map{
			"error": "cannot create todo",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(

		fiber.Map{
			"message": "Todo created",
			"todo":    todo,
		},
	)
}
func GetTodos(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)

	cursor, err := db.DB.Collection("todos").Find(c.Context(), bson.M{
		"userId": userId,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch todos",
		})
	}
	var todos []bson.M
	if err := cursor.All(c.Context(), &todos); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "cannot parse todos",
		})
	}
	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"todos": todos,
		},
	)
}
func DeleteTodo(c *fiber.Ctx) error {
	todoId := c.Params("id")
	userId := c.Locals("userId").(string)
	objId, err := primitive.ObjectIDFromHex(todoId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid todo id",
		})
	}
	filter := bson.M{
		"_id":    objId,
		"userId": userId,
	}
	result, err := db.DB.Collection("todos").DeleteOne(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "cannot delete todo",
		})
	}
	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "todo not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "todo deleted",
	})
}
