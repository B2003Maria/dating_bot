package main

import (
  // "log"
  "github.com/go-telegram-bot-api/telegram-bot-api/v5"
  "fmt"
)

func (b *Bot) handleMessage(update tgbotapi.Update) {
    chatID := update.Message.Chat.ID

    // Проверяем, если пользователь уже в процессе заполнения анкеты
    if state, exists := userStates[chatID]; exists && state.CurrentField != "" {
        b.handleProfileResponse(update)
        return
    }

    // Обработка команды /start и кнопки "Добавить анкету"
    switch update.Message.Text {
    case "/start":
        b.handleStart(update)
    case "Добавить анкету":
        b.handleProfileWriting(update)
    }
}

// Реакция на команду /start
func (b *Bot) handleStart(update tgbotapi.Update) {
    text := fmt.Sprintf("Привет, %s! В данном боте ты можешь заводить новые знакомства!", update.SentFrom().FirstName)

    // Создаем кнопку "Добавить анкету"
    keyboard := tgbotapi.NewReplyKeyboard(
        tgbotapi.NewKeyboardButtonRow(
            tgbotapi.NewKeyboardButton("Добавить анкету"),
        ),
    )

    msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
    msg.ReplyMarkup = keyboard
    b.api.Send(msg)
}

// Хранение информации о состоянии пользователя при заполнении анкеты
type UserState struct {
  Profile      Profile
  CurrentField string
}

// Хранение состояний пользователей
var userStates = make(map[int64]*UserState)

func (b *Bot) handleProfileWriting(update tgbotapi.Update) {
    chatID := update.Message.Chat.ID

    // Создаем новое состояние для пользователя, если его нет
    if _, exists := userStates[chatID]; !exists {
        userStates[chatID] = &UserState{
            Profile:      Profile{},
            CurrentField: "Name",
        }
    }

    // Приветственное сообщение
    msg := tgbotapi.NewMessage(chatID, "Отлично, давай заполним информацию о тебе, которая будет отображаться в анкете.")
    b.api.Send(msg)
    b.askNextQuestion(chatID)
}

// Функция для запроса следующего вопроса
func (b *Bot) askNextQuestion(chatID int64) {
    state := userStates[chatID]

    // Определение вопроса и клавиатуры в зависимости от текущего поля
    var question string
    var replyMarkup interface{}
    
    switch state.CurrentField {
    case "Name":
        question = "Как тебя зовут?"
    case "Age":
        question = "Сколько тебе лет?\nВведите число от 18 до 99."
    case "Gender":
        question = "Укажи свой пол?"
        // Создаем клавиатуру с кнопками для выбора пола
        replyMarkup = tgbotapi.NewReplyKeyboard(
            tgbotapi.NewKeyboardButtonRow(
                tgbotapi.NewKeyboardButton("Мужской"),
                tgbotapi.NewKeyboardButton("Женский"),
            ),
        )
    case "Interest":
        question = "Кто тебе интересен?"
        // Создаем клавиатуру с кнопками для выбора интересов
        replyMarkup = tgbotapi.NewReplyKeyboard(
            tgbotapi.NewKeyboardButtonRow(
                tgbotapi.NewKeyboardButton("Девушки"),
                tgbotapi.NewKeyboardButton("Парни"),
                tgbotapi.NewKeyboardButton("Неважно"),
            ),
        )
    case "Description":
        question = "Напиши краткое описание о себе"
    case "Photo":
        question = "Загрузите своё фото (Пока в разработке...Отправье любое сообщение)"
    default:
        question = "Анкета заполнена! Спасибо."
        b.finishProfile(chatID)
        return
    }

    // Отправка вопроса с клавиатурой (если она есть)
    msg := tgbotapi.NewMessage(chatID, question)
    if replyMarkup != nil {
        msg.ReplyMarkup = replyMarkup
    }
    b.api.Send(msg)
}

func (b *Bot) handleProfileResponse(update tgbotapi.Update) {
    chatID := update.Message.Chat.ID
    state := userStates[chatID]

    // Сохранение ответа в зависимости от текущего поля и проверка на корректность
    switch state.CurrentField {
    case "Name":
        state.Profile.Name = update.Message.Text
        state.CurrentField = "Age"
    case "Age":
        var age uint8
        _, err := fmt.Sscanf(update.Message.Text, "%d", &age)
        if err != nil || age < 18 || age > 99 {
            msg := tgbotapi.NewMessage(chatID, "Пожалуйста, введи корректный возраст (от 18 до 99).")
            b.api.Send(msg)
            return
        }
        state.Profile.Age = age
        state.CurrentField = "Gender"
    case "Gender":
        if update.Message.Text != "Мужской" && update.Message.Text != "Женский" {
            msg := tgbotapi.NewMessage(chatID, "Пожалуйста, выбери пол из предложенных вариантов: Мужской или Женский.")
            b.api.Send(msg)
            return
        }
        state.Profile.Gender = update.Message.Text
        state.CurrentField = "Interest"
    case "Interest":
        if update.Message.Text != "Девушки" && update.Message.Text != "Парни" && update.Message.Text != "Неважно" {
            msg := tgbotapi.NewMessage(chatID, "Пожалуйста, выбери один из предложенных вариантов: Девушки, Парни, Неважно.")
            b.api.Send(msg)
            return
        }
        state.Profile.Interest = update.Message.Text
        state.CurrentField = "Description"
    case "Description":
        state.Profile.Description = update.Message.Text
        state.CurrentField = "Photo"
    case "Photo":
        // msg := tgbotapi.NewMessage(chatID, "Пока в разработке. Отправь любое сообщение")
        // b.api.Send(msg)
        state.Profile.Photo = "Пока в разработке..."
        state.CurrentField = ""
    }

    b.askNextQuestion(chatID)
}

// Завершение заполнения анкеты
func (b *Bot) finishProfile(chatID int64) {
    state := userStates[chatID]
    profile := state.Profile

    summary := fmt.Sprintf(
        "Анкета заполнена!\nИмя: %s\nВозраст: %d\nПол: %s\nИнтересы: %s\nОписание: %s\nФото: %s",
        profile.Name, profile.Age, profile.Gender, profile.Interest, profile.Description, profile.Photo,
    )

    msg := tgbotapi.NewMessage(chatID, summary)
    b.api.Send(msg)
    msg = tgbotapi.NewMessage(chatID, "Дальше будет добавлен функционал изменения анкеты")
    b.api.Send(msg)
  // replyMarkup := tgbotapi.NewReplyKeyboard(
  //           tgbotapi.NewKeyboardButtonRow(
  //               tgbotapi.NewKeyboardButton("Да"),
  //               tgbotapi.NewKeyboardButton("Исправить"),
  //               tgbotapi.NewKeyboardButton("Неважно"),
  //           ),
  //       )
    // Удаление состояния пользователя после завершения заполнения анкеты
    delete(userStates, chatID)
}
