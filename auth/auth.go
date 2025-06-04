package auth

import (
	"golang.org/x/crypto/bcrypt"
)

var contextKey = "user"


// Check if two passwords match using Bcrypt's CompareHashAndPassword
func DoPasswordsMatch(hashedPassword, currPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword), []byte(currPassword))
	return err == nil
}