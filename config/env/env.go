package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func Load() {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func getKey(key string, fallback string) string {
	value, ok := os.LookupEnv(key)

	if !ok {
		return fallback
	}

	return value
}

func GetString(key string, fallback string) string {
	return getKey(key, fallback)
}

func GetInt(key string, fallback int) int {
	value := getKey(key, strconv.Itoa(fallback))

	intValue, err := strconv.Atoi(value)

	if err != nil {
		fmt.Printf("Error converting %s to int: %v\n", key, err)
		return fallback
	}

	return intValue
}

func GetBool(key string, fallback bool) bool {
	value := getKey(key, strconv.FormatBool(fallback))

	boolValue, err := strconv.ParseBool(value)

	if err != nil {
		fmt.Printf("Error converting %s to bool: %v\n", key, err)
		return fallback
	}

	return boolValue
}