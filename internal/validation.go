package internal

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func isUsernameValid(s string) bool {
	if len(RemoveSpace(s)) == 0 || len(s) > 20 {
		return false
	}
	loginConvention := "^[a-zA-Z0-9]*[-]?[a-zA-Z0-9]*$"
	if re, _ := regexp.Compile(loginConvention); !re.MatchString(s) {
		return false
	}
	return true
}

func RemoveSpace(s string) string {
	return strings.TrimSpace(s)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", fmt.Errorf("HashPassword: %w", err)
	}
	return string(bytes), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
