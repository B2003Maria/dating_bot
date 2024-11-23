package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api *tgbotapi.BotAPI
}

func main() {
	initMongo()

	user := &User{
		UserID:   12345,
		ChatID:   67890,
		Username: "test_user",
		Profile: Profile{
			Name:        "John",
			Age:         30,
			Gender:      "Male",
			Interest:    "Reading",
			Description: "Avid reader and traveler",
			Photo:       "photo_id_123",
		},
		LikedTo:   make(map[int64]struct{}),
		LikedBy:   make(map[int64]struct{}),
		DislikeTo: make(map[int64]struct{}),
		SleptTo:   make(map[int64]struct{}),
	}

	if err := CreateUser(user); err != nil {
		lg.Fatalf("Failed to create user! %s", err)
	}
	bot := initBot(token)
	bot.start()
}

// Инициализация бота
func initBot(token string) *Bot {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		lg.Panicf("Failed to read token: %s", err)
	}
	api.Debug = true
	log.Printf("Авторизован под аккаунтом %s", api.Self.UserName)
	return &Bot{api: api}
}

// Запуск обработки обновлений
func (b *Bot) start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			b.handleCallbackQuery(update)
		}
		if update.Message != nil {
			b.handleMessage(update)
		}
	}
}
