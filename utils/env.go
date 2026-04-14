package utils

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from a .env file in the current directory.
// Uses godotenv to parse the file and set environment variables.
// Fatally exits if the .env file cannot be loaded.
// Note: This function is currently not used in the main application flow.
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
