package main

import (
	"log"
	"os"

	"github.com/kamva/mgm/v3"
)

// Логгер для вывода
var lg *log.Logger = log.New(os.Stdout, "INFO: ", log.Ltime|log.Lmicroseconds)

// Profile структура для профиля пользователя
type Profile struct {
	Name        string `bson:"name"`
	Age         uint8  `bson:"age"`
	Gender      string `bson:"gender"`
	Interest    string `bson:"interest"`
	Description string `bson:"description"`
	Photo       string `bson:"photo"`
}

// User структура для пользователя
type User struct {
	// Встроенные базовые поля (ID, CreatedAt, UpdatedAt)
	mgm.DefaultModel `bson:",inline"`

	UserID    int64              `bson:"user_id"`
	ChatID    int64              `bson:"chat_id"`
	Username  string             `bson:"username"`
	Profile   Profile            `bson:"profile"`
	LikedTo   map[int64]struct{} `bson:"liked_to"`
	LikedBy   map[int64]struct{} `bson:"liked_by"`
	DislikeTo map[int64]struct{} `bson:"dislike_to"`
	SleptTo   map[int64]struct{} `bson:"slept_to"`
}

// Users массив пользователей
var Users []User
