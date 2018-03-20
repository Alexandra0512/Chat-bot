package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

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

//структура json
type jsonStruct struct {
	keyWord    string `json:"keyWord"`    //ключевое слово из сообщения
	amount     string `json:"amount"`     //сумма, указанная в сообщении
	targetName string `json:"targetName"` //имя цели, если ключевое слово-цель
	date       string `json:"date"`       //дата, указанная в сообщении
}

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

func searchKey(sms string) {
	var (
		key     string
		key_id  int
		mas_sms []string
		mas_key []string
		tmp     []string
	)

	key = "справка spravka зачисление zachislenie пополнение popolnenie поступление postuplenie списание Spisanie покупка Perevod Оплата Oplata Цель tsel Баланс balans Остаток Ostatok Выписка Vypiska Редектировать Redektirovat Изменить Izmenit Кошелек Koshelek Долг Dolg Создать Sozdat Выход vykhod"
	//sms= "Schet *1096: postuplenie zarabotnoy plati 5605.25 RUB; 02.03.2018 05:36:31; Dostupno 5938.37 RUB. Detali v mobilnom banke: vtb.ru/app"
	//sms= "VISA6679 23.02.18 04:38 списание 7р KOPILKA KARTS-VKLAD Баланс: 3973р"
	//sms= "VISA6679 23.02.18 04:38 списание 7р KOPILKA KARTS-VKLAD Баланс: 3973р"

	mas_sms = strings.Split(sms, " ") //переход от строки к массиву
	mas_key = strings.Split(key, " ")
	//fmt.Print(len(mas_key))
	for i := 0; i < len(mas_sms); i++ {
		for j := 0; j < len(mas_key); j++ {
			//fmt.Println("mas_sms[i]=", mas_sms[i])
			//fmt.Print("mas_key[ij=", mas_key[j])
			if mas_sms[i] == mas_key[j] {
				//fmt.Print("j=", j)
				key_id = j
				tmp = append(tmp, mas_sms[i])
				if (len(tmp)) >= 2 {
					fmt.Printf("%q\n", "Некорректный запрос")
					break
				} else {
					fmt.Printf("%q\n", mas_sms[i])
				}
			}

		}
	}
	//fmt.Print("jyfuj  ", key_id)
	switch key_id {
	case 2, 3, 4, 5, 6, 7:
		structSMS(mas_sms)
		break
	}
}

func structSMS(mas_sms []string) {
	var (
		mas_result []string
	)
	//fmt.Print("jsdssyfuj")
	if strings.Contains(mas_sms[0], "VISA") {

		data := mas_sms[1]
		time := mas_sms[2]
		sum := strings.Replace(mas_sms[4], "р", "", -1)
		tipe := "undefiend"
		mas_result = append(mas_result, tipe, data, time, sum)
	} else if strings.Contains(mas_sms[0], "Schet") {
		tipe := mas_sms[3]
		data := mas_sms[7]
		time := mas_sms[8]
		sum := mas_sms[5]
		mas_result = append(mas_result, tipe, data, time, sum)
	} else if strings.Contains(mas_sms[0], "MIR") {
		tipe := mas_sms[3]
		data := mas_sms[1]

		time := mas_sms[2]
		sum := strings.Replace(mas_sms[4], "р", "", -1)
		mas_result = append(mas_result, tipe, data, time, sum)
	}

	//fmt.Println("", mas_result)
	//return
	convertMassToJson(mas_result)
}

func convertMassToJson(array []string) {
	right := &jsonStruct{
		keyWord:    array[0],
		amount:     array[1],
		targetName: array[2],
		date:       array[3],
	}
	fmt.Println("", right)
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
			break
		case "start":
		case "file":
		}

		var sms string
		sms = "VISA6679 23.02.18 04:38 зачисление 3000р Баланс: 4254.87р"
		searchKey(sms)
		// обработка сообщений

		// парсинг ключевых слов

		// вызов определенных функций

	}
}
