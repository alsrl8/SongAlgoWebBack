package utils

import "os"

func IsDevelopmentMode() bool {
	env := GetEnvWithDefault("ENV_TYPE", "DEV")
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
