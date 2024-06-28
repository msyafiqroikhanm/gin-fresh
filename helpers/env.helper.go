package helpers

import "os"

// Function to abstarct the process of getting env data that also set default value
// if the env attribute empty or doesn't exist
func GetENVWithDefault(key string, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}
