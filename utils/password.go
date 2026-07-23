package utils

import (
	"github.com/gofiber/fiber/v3/log"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain text password using bcrypt.
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// Jangan pakai log.Fatalf: itu mematikan seluruh proses server.
		// bcrypt bisa gagal (mis. password > 72 byte); kembalikan error saja.
		log.Errorf("Error hashing password: %v", err)
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPassword compares a hashed password with a plain text password.
func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Errorf("Error comparing passwords: %v", err)
	}
	return err == nil
}
