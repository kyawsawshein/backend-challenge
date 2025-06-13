package user

import (
	"errors"
	"fmt"
	"time"

	"backend-challenge/auth"
	"backend-challenge/db"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

const USER = "users"

func SetUpRoutes(webApp *fiber.App) {
	userGroup := webApp.Group("/user")
	userGroup.Post("/register", register)
	userGroup.Get("/", listUsers)
	userGroup.Get("/:id", getUser)
	userGroup.Put("/:id", updateUser)
	userGroup.Delete("/:id", deleteUser)
}

func GetUserCollection() *mongo.Collection {
	return db.Connection.Collection(USER)
}

func GetUser(c *fiber.Ctx, email string) (*User, error) {
	var user User
	userCollection := db.GetCollection(USER)
	if err := userCollection.FindOne(c.Context(), bson.M{"email": email}).Decode(&user); err != nil {
		return nil, c.Status(500).JSON(fiber.Map{"error": "User not found", "detail": err})
	}
	return &user, nil
}

func register(c *fiber.Ctx) error {
	ctx := c.Context()
	var payload struct {
		Name            string `json:"name"`
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input!"})
	}
	if payload.Password != payload.ConfirmPassword {
		return c.Status(300).JSON(fiber.Map{"errro": "Password does not match!"})
	}
	count, _ := GetUserCollection().CountDocuments(ctx, bson.M{"email": payload.Email})
	if count > 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Email already registered~"})
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	user := User{
		Name:      payload.Name,
		Email:     payload.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}
	_, err := GetUserCollection().InsertOne(ctx, user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "User create filled!", "detail": err})
	}
	return c.JSON(fiber.Map{"message": "User created"})
}

func listUsers(c *fiber.Ctx) error {
	user_id, err := auth.ValidateJwt(c)
	if err != nil  {
		return c.Status(422).JSON(fiber.Map{"Unauthal": err})
	}
	// TODO user role
	fmt.Println("User Id : ", user_id)
	ctx := c.Context()
	cursor, err := GetUserCollection().Find(ctx, bson.M{})
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "User not found", "detail": err})
	}
	defer cursor.Close(ctx)

	var users []User
	for cursor.Next(ctx) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			continue
		}
		user.Password = ""
		users = append(users, user)
	}
	return c.JSON(users)
}

func getUser(c *fiber.Ctx) error {
	ctx := c.Context()
	user_id, err := auth.ValidateJwt(c)
	if err != nil  {
		return c.Status(422).JSON(fiber.Map{"Unauthal": err})
	}
	// TODO user role
	fmt.Println("User Id : ", user_id)
	userID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID", "detail": err})
	}
	cursor, err := GetUserCollection().Find(ctx, bson.M{"_id": userID})
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "User not found", "detail": err})
	}
	defer cursor.Close(ctx)

	var user User
	for cursor.Next(ctx) {
		if err := cursor.Decode(&user); err != nil {
			errors.New(fmt.Sprintf("Error get user."))
		}
		user.Password = ""
	}
	return c.JSON(user)
}

func updateUser(c *fiber.Ctx) error {
	ctx := c.Context()
	user_id, err := auth.ValidateJwt(c)
	if err != nil  {
		return c.Status(422).JSON(fiber.Map{"Unauthal": err})
	}
	// TODO user role
	fmt.Println("User Id : ", user_id)
	userID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}
	var payload User
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"})
	}
	updateDoc := bson.M{"$set": StructToBsonM(payload)}
	if updateDoc["password"] != nil {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
		updateDoc["password"] = string(hashedPassword)
	}
	if len(updateDoc) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "No fields to update?"})
	}
	_, err = GetUserCollection().UpdateOne(
		ctx,
		bson.M{"_id": userID},
		updateDoc,
	)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Updaate failed", "detail": err})
	}
	return c.JSON(fiber.Map{"message": "User updated."})
}

func deleteUser(c *fiber.Ctx) error {
	ctx := c.Context()
	user_id, err := auth.ValidateJwt(c)
	if err != nil  {
		return c.Status(422).JSON(fiber.Map{"Unauthal": err})
	}
	// TODO user role
	fmt.Println("User Id : ", user_id)
	userId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID", "detail": err})
	}
	res, err := GetUserCollection().DeleteOne(ctx, bson.M{"_id": userId})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Delete failed", "detail": err})
	}
	if res.DeletedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "User not found."})
	}
	return c.JSON(fiber.Map{"message": "User deleted"})
}
