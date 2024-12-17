package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api *tgbotapi.BotAPI
}

func main() {
	loadEnv()
	err := loadFromFile(&Users)
	if err != nil {
		lg.Fatalf("Failed to load db: %s", err)
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
	log.Printf("Authorised under account %s", api.Self.UserName)
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
