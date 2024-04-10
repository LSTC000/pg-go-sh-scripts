package server

import (
	"fmt"

	"github.com/joho/godotenv"
)

func setDotEnv() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("cannot load .env file: %w", err)
	}
	return nil
}
