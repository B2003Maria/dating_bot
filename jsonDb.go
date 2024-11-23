package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func saveToFile(data map[int64]User) error {
	// Сериализация данных в JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка сериализации: %v", err)
	}

	err = os.WriteFile(dbPath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("ошибка записи в файл: %v", err)
	}
	return nil
}

func loadFromFile(data *map[int64]User) error {
	// Чтение файла
	jsonData, err := os.ReadFile(dbPath)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла: %v", err)
	}

	err = json.Unmarshal(jsonData, data)
	if err != nil {
		return fmt.Errorf("ошибка десериализации: %v", err)
	}
	return nil
}
