package service

import (
	"backend-challenge/db"
	"backend-challenge/user"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUser(c *fiber.Ctx, email string) (*user.User, error) {
	var user user.User
	userCollection := db.GetCollection("users")
	if err := userCollection.FindOne(c.Context(), bson.M{"email": email}).Decode(&user); err != nil {
		return nil, c.Status(500).JSON(fiber.Map{"error": "User not found", "detail": err})
	}
	return &user, nil
}