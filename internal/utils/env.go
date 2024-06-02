package utils

import (
	"os"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func GetEnv(key string) string {
	value := os.Getenv(key)

	if value == "" {
		panic("Missing " + key)
	}

	return value
}
