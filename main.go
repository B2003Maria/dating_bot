package main

import (
    "log"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
    api *tgbotapi.BotAPI
}

func main() {
    bot := initBot(token)
    bot.start()
}

// Инициализация бота
func initBot(token string) *Bot {
    api, err := tgbotapi.NewBotAPI(token)
    if err != nil {
        log.Panic(err)
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
        if update.Message != nil {
            b.handleMessage(update)
        }
    }
}
