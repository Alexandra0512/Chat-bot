package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	drive "google.golang.org/api/drive/v3"
	"gopkg.in/telegram-bot-api.v4"
)

var (
	bot        *tgbotapi.BotAPI
	srvDrive   *drive.Service
	ChatID     int64
	UzverID    int
	isReadAuth bool
	updates    tgbotapi.UpdatesChannel
)

// SendMessageToTelegram вывод сообщения в телеграмме от бота
func SendMessageToTelegram(textMessage string) {

	// формирование сообщения
	msg := tgbotapi.NewMessage(ChatID, textMessage)

	// вывод сообщения в телеграмме
	bot.Send(msg)

}

// initBot Инициализация бота
func initBot() {
	// Считывание файла с настройками подключения к Telegram
	file, _ := os.Open("config_telegram.json")
	decoder := json.NewDecoder(file)

	// Считывание данных из файла настроек
	configuration := Config{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Panic(err)
	}

	// создание экземпляра класса чат-бота
	bot, err = tgbotapi.NewBotAPI(configuration.TelegramBotToken)
	if err != nil {
		log.Panic(err)
	}
}

func main() {

	isReadAuth = false
	initBot()

	bot.Debug = false
	log.Printf("Авторизация в аккаунте: %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Panic(err)
	}

	// В канал updates будут приходить все новые сообщения.
	for update := range updates {
		// Создав структуру - можно её отправить обратно боту

		fmt.Printf("%s\n", update.Message.From.ID)

		if update.Message == nil {
			continue
		}

		if isReadAuth {
			code := update.Message.Text
			isReadAuth = false
			AuthCode <- code
		}

		// обработка команд
		// P.S. команда начинается с "/"
		switch update.Message.Command() {
		case "auth": // авторизация пользователя
			ChatID = update.Message.Chat.ID
			UzverID = update.Message.From.ID
			go initGoogle()
			isReadAuth = true
		case "start":
		case "file":
		}

		// обработка сообщений

		// парсинг ключевых слов

		// вызов определенных функций

	}
}
