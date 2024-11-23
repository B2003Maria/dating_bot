package main

import (
	"os"

	"github.com/joho/godotenv"
)

func checkAge(age uint8) bool {
	if age >= 18 {
		return true
	}
	return false
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		lg.Fatalf("Failed to load .env file: %v", err)
	}

	// Инициализируем глобальные переменные
	token = os.Getenv("TOKEN")
	dbPath = os.Getenv("DBPATH")
}
