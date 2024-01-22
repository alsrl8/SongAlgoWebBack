package utils

import (
	"errors"
	"fmt"
	"os"
)

func IsDevelopmentMode() bool {
	env := GetEnvWithDefault("ENV_TYPE", "DEV")
	fmt.Printf("ENV_TYPE: %s\n", env)
	if env == "DEV" {
		return true
	} else {
		return false
	}
}

func GetEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", errors.New("environment variable not set")
	}
	return value, nil
}
