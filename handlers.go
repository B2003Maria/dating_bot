package main

import (
	// "log"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleMessage(update tgbotapi.Update) {
	if state, exists := userStates[update.Message.From.ID]; exists && state.CurrentField != "AllGood" {
		b.handleProfileResponse(update)
		return
	}
	// Обработка команды /start и кнопки "Добавить анкету"
	switch update.Message.Text {
	case "/start":
		if _, exists := Users[update.Message.From.ID]; exists {
			fmt.Printf("Пользователь с ID=%d существует.\n", update.Message.From.ID)
			b.handleExistingStart(update)
			return
		}
		Users[update.Message.From.ID] = User{}
		saveToFile(Users)
		b.handleFirstStart(update)
	case "Добавить анкету":
		b.handleProfileWriting(update)
	case "Смотреть анкеты":
		b.handleWatchProfile(update)
	}
}

func (b *Bot) handleCallbackQuery(update tgbotapi.Update) {
	callback := update.CallbackQuery

	// Подтверждаем получение CallbackQuery, чтобы Telegram перестал показывать "часики"
	ack := tgbotapi.NewCallback(callback.ID, "") // Пустая строка, если не нужно отправлять сообщение пользователю
	if _, err := b.api.Request(ack); err != nil {
		lg.Printf("Ошибка подтверждения CallbackQuery: %v", err)
		return
	}

	// Обрабатываем действие на основе данных кнопки
	action := callback.Data
	lg.Println(action)
	var response string

	switch action {
	case "like":
		response = "Вы поставили ❤️"
	case "dislike":
		response = "Вы поставили 💔"
	case "sleep":
		response = "Вы выбрали 😴. Пользователь будет отложен."
	default:
		response = "Неизвестное действие."
	}

	// Отправляем ответ пользователю
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, response)
	b.api.Send(msg)
}

func (b *Bot) showProfile(update tgbotapi.Update, profile Profile) {
	summary := fmt.Sprintf(
		"Имя: %s\nВозраст: %d\nПол: %s\nИнтересы: %s\nОписание: %s",
		profile.Name, profile.Age, profile.Gender, profile.Interest, profile.Description,
	)

	// Отправка фотографии с подписью
	photoMsg := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FileID(profile.Photo))
	photoMsg.ReplyMarkup = ReactionsInlineKeyboard
	photoMsg.Caption = summary
	b.api.Send(photoMsg)
}

func (b *Bot) showProfiles(update tgbotapi.Update) {
	for _, user := range Users {
		lg.Println(user)
		b.showProfile(update, user.Profile)
	}
}

// Реакция на команду /start
func (b *Bot) handleFirstStart(update tgbotapi.Update) {
	text := fmt.Sprintf("Привет, %s! В данном боте ты можешь заводить новые знакомства!", update.SentFrom().FirstName)

	// Создаем кнопку "Добавить анкету"
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ReplyMarkup = AddProfileKeyboard
	b.api.Send(msg)
}

func (b *Bot) handleExistingStart(update tgbotapi.Update) {
	text := fmt.Sprintf("Привет снова, %s! В данном боте ты можешь заводить новые знакомства!", update.SentFrom().FirstName)

	// Создаем кнопку "Добавить анкету"
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ReplyMarkup = WatchProfilesKeyboard
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
	msg := tgbotapi.NewMessage(chatID, "Отлично, давай заполним информацию о тебе, которая будет отображаться в анкете.")
	b.api.Send(msg)
	b.askNextQuestion(update)
}

// Функция для запроса следующего вопроса
func (b *Bot) askNextQuestion(update tgbotapi.Update) {
	state := userStates[update.Message.Chat.ID]

	// Определение вопроса и клавиатуры в зависимости от текущего поля
	var question string
	var keyboard interface{}

	switch state.CurrentField {
	case "Name":
		question = "Как тебя зовут?"
	case "Age":
		question = "Сколько тебе лет?\nВведите число от 18 до 99."
	case "Gender":
		question = "Укажи свой пол?"
		keyboard = ChooseGenderKeyboard
	case "Interest":
		question = "Кто тебе интересен?"
		keyboard = ChooseInterestKeyboard
	case "Description":
		question = "Напиши краткое описание о себе"
	case "Photo":
		question = "Загрузи своё фото"
	case "Check":
		question = "Всё введено верно?"
		keyboard = IsProfileOkKeyboard
	case "AllGood":
		question = "Спасибо! Ваша анкета сохранена!"
		keyboard = WatchProfilesKeyboard
		b.saveUser(update, state.Profile)
	default:
		b.finishProfile(update)
		return
	}
	// Отправка вопроса с клавиатурой (если она есть)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, question)
	if keyboard != nil {
		msg.ReplyMarkup = keyboard
	} else {
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	}
	b.api.Send(msg)
}

func (b *Bot) handleProfileResponse(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	state := userStates[chatID]
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
		if update.Message.Photo != nil {
			state.Profile.Photo = update.Message.Photo[0].FileID
			state.CurrentField = "" // Заполнение анкеты завершено
		} else {
			msg := tgbotapi.NewMessage(chatID, "Пожалуйста, загрузи фото.")
			b.api.Send(msg)
			return
		}
	case "Check":
		// lg.Println("\n\n\n\n")
		if update.Message.Text == "Всё верно" {
			state.CurrentField = "AllGood"
			b.askNextQuestion(update)
			return
		}
		if update.Message.Text == "Ввести анекту заново" {
			state.CurrentField = "Name"
			b.askNextQuestion(update) // Перезапускаем вопросы
			return
		}
		// Если текст не распознан
		msg := tgbotapi.NewMessage(chatID, "Пожалуйста, выбери один из вариантов: 'Всё верно' или 'Ввести анкету заново'.")
		b.api.Send(msg)
	case "AllGood":
		b.saveUser(update, state.Profile)
		b.handleMessage(update)
		return
	}
	b.askNextQuestion(update)
}

// Завершение заполнения анкеты
func (b *Bot) finishProfile(update tgbotapi.Update) {
	state := userStates[update.Message.Chat.ID]
	profile := state.Profile

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Анкета заполнена! Спасибо!\nДавай сверим информацию!")
	b.api.Send(msg)
	summary := fmt.Sprintf(
		"Имя: %s\nВозраст: %d\nПол: %s\nИнтересы: %s\nОписание: %s",
		profile.Name, profile.Age, profile.Gender, profile.Interest, profile.Description,
	)

	// Отправка фотографии с подписью
	photoMsg := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FileID(profile.Photo))
	photoMsg.Caption = summary
	b.api.Send(photoMsg)

	// Если анкета подтверждена, удаляем состояние
	if state.CurrentField == "AllGood" {
		b.saveUser(update, profile)
	} else {
		state.CurrentField = "Check"
		b.askNextQuestion(update)
	}
}

func (b *Bot) handleWatchProfile(update tgbotapi.Update) {
	lg.Println(Users)
	b.showProfiles(update)
}

func (b *Bot) saveUser(update tgbotapi.Update, profile Profile) {
	var user User
	user.ChatID = update.Message.Chat.ID
	user.Profile = profile
	user.Username = update.Message.From.UserName
	Users[update.Message.From.ID] = user
	err := saveToFile(Users)
	if err != nil {
		lg.Printf("Failed to save user %v:\n%s", user, err)
	}
	lg.Printf("User: %s is saved!", update.Message.From.UserName)
	delete(userStates, update.Message.Chat.ID)
}
