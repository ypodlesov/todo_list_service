package base

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("1.com", "http://1.com"),
		tgbotapi.NewInlineKeyboardButtonData("2", "2"),
		tgbotapi.NewInlineKeyboardButtonData("3", "3"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("4", "4"),
		tgbotapi.NewInlineKeyboardButtonData("5", "5"),
		tgbotapi.NewInlineKeyboardButtonData("6", "6"),
	),
)

func processUpdate(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	// Check if we've gotten a message update.
	var msg tgbotapi.MessageConfig
	if update.Message != nil {
		// Construct a new message from the given chat ID and containing
		// the text that we received.
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		// If the message was open, add a copy of our numeric keyboard.
		switch update.Message.Text {
		case "open":
			msg.ReplyMarkup = numericKeyboard

		}
	} else if update.CallbackQuery != nil {
		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
		if _, err := bot.Request(callback); err != nil {
			panic(err)
		}
		msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
	}
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}

func Run() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	// bot.ListenForWebhook()

	// Loop through each update.
	for update := range updates {
		go processUpdate(bot, &update)
	}
}
