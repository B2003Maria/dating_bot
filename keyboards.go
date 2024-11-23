package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var ReactionsInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("❤️", "like"),
		tgbotapi.NewInlineKeyboardButtonData("💔", "dislike"),
		tgbotapi.NewInlineKeyboardButtonData("😴", "sleep"),
	),
)

var AddProfileKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Добавить анкету"),
	),
)

var ChooseGenderKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Мужской"),
		tgbotapi.NewKeyboardButton("Женский"),
	),
)

var ChooseInterestKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Девушки"),
		tgbotapi.NewKeyboardButton("Парни"),
		tgbotapi.NewKeyboardButton("Неважно"),
	),
)

var IsProfileOkKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Всё верно"),
		tgbotapi.NewKeyboardButton("Ввести анекту заново"),
	),
)

var WatchProfilesKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Смотреть анкеты"),
	),
)
