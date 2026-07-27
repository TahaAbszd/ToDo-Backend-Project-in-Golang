package controller

import (
	"todo/src/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateTodo(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)
	type body struct{
		Title string `json:"title"`
		Description string `json:"description"`
		Status string `json:"status"`
	}
	var data body
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":"cannot parse json",
		})
	}
	if data.Status != string(models.StatusCompleted) && 
		data.Status != string(models.StatusINcomplete) {
			data.Status = string(models.StatusINcomplete)
		}
	todo := bson.M{
		"_id:" primitive.NewObjectID(),
		"title:" data.Title,
		"description:" data.Description,
		"status:" data.Status,
		"userId:" userId,
		}

	_,err := db.DB.Collection("todos").InsertOne(c.Context(),todo)
	if err != nil {
		return c.Status(fiber.StatusInsufficientStorage).JSON(fiber.Map{
			"error":"cannot create todo"
		})
	}
	return c.Status(fiber.StatusCreated).JSON(

		fiber,Map{
			"message":"Todo created",
			"todo":todo,
		},
	)
}
func GetTodos(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)

	cursor , err := db.DB.Collection("todos").Find(c.Context(),bson.M{
		"userId":userId,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch todos",
		})
	}
	var todos []bson.M
	if err :=cursor.All(c.Context(),&todos); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":"cannot parse todos"
		})
	}
	return c.Status(fiber.StatusOk).JSON(
		fiber.Map{
			"todos":todos,
		},
	)
}