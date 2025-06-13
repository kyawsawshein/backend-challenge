package auth

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var contextKey = "user"


// Check if two passwords match using Bcrypt's CompareHashAndPassword
func DoPasswordsMatch(hashedPassword, currPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword), []byte(currPassword))
	return err == nil
}


func ValidateJwt(c *fiber.Ctx) (string, error) {
	fmt.Println("token : ", c.Locals("token"))
	userToken, ok := c.Locals("token").(*jwt.Token)
	if !ok {
		return "", fiber.ErrUnauthorized
	}
	
	if userToken != nil {
		claims := userToken.Claims.(jwt.MapClaims)
		tokenUserId := claims["user_id"].(string)
			return tokenUserId, nil
		}
	return "", c.Status(500).JSON(fiber.Map{"error": "Token expiry or empty"})
}
