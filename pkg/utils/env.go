package utils

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

func CheckAdminRole(key string) (string, error) {
	err := godotenv.Load(".env")

	if err != nil {
		return "", errors.New("Error loading .env file")
	}
  
	return os.Getenv(key), nil
}