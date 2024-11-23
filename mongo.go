package main

import (
	"github.com/kamva/mgm/v3"                   // mgm фреймворк
	"go.mongodb.org/mongo-driver/mongo/options" // настройки клиента
)

func initMongo() {
	// Инициализация mgm с указанием подключения к MongoDB
	err := mgm.SetDefaultConfig(nil, "dating_app", options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		lg.Fatalf("Ошибка подключения к MongoDB: %v", err)
	}
}

func CreateUser(user *User) error {
	return mgm.Coll(user).Create(user)
}
