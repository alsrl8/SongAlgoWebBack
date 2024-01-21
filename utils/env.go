package utils

import "os"

func IsDevelopmentMode() bool {
	env := GetEnvWithDefault("env-type", "dev")
	if env == "dev" {
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
