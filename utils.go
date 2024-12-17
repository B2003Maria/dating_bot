package main

import (
	"os"

	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		lg.Fatalf("Failed to load .env file: %v", err)
	}

	// Инициализируем глобальные переменные
	token = os.Getenv("TOKEN")
	dbPath = os.Getenv("DBPATH")
}
