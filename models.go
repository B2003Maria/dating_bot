package main

import (
	"log"
	"os"
)

// Логгер для вывода
var lg *log.Logger = log.New(os.Stdout, "INFO: ", log.Ltime|log.Lmicroseconds)

var (
	dbPath string
	token  string
)

// Profile структура для профиля пользователя
type Profile struct {
	Name        string
	Age         uint8
	Gender      string
	Interest    string
	Description string
	Photo       string
}

// User структура для пользователя
type User struct {
	ChatID    int64
	Username  string
	Profile   Profile
	LikedTo   map[int64]struct{}
	LikedBy   map[int64]struct{}
	DislikeTo map[int64]struct{}
	SleptTo   map[int64]struct{}
}

// Users массив пользователей
var Users = make(map[int64]User)
