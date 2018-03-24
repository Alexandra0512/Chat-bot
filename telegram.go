package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	drive "google.golang.org/api/drive/v3"
	sheets "google.golang.org/api/sheets/v4"
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

	bot.Debug = true

	// _, err = bot.SetWebhook(tgbotapi.NewWebhookWithCert("https://www.google.com:8443/"+bot.Token, "cert.pem"))
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

// ParseMess парсинг сообщения от банка.
// На вход сообщения от банка
// на выходе данные в формате jsonStruct
func ParseMess(message string) (data jsonStruct) {

	fmt.Println("В процедуру ParseMees пришло: " + message)

	regexpSeparator := " "
	// тип карты
	regexpCard := "((MIR|VISA|MAES|ECMC)-?\\d{4})"
	// дата совершения операции
	regexpDate := "((0[1-9]|[12][0-9]|3[01]).(0[1-9]|1[012]).(\\d{2}))"
	// время совершения операции
	regexpTime := "((0[0-9]|1[0-9]|2[0-3]]):([0-5][1-9]))"
	// опреация
	regexpKeyWord := "((списание|покупка|зачисление зарплаты|зачисление пенсии|зачисление|поступление|выдача наличных|выдача|возврат покупки))"
	// сумма операции
	regexpSum := "((\\d+(.\\d{2})?)+р)"
	//
	regexpW := "(([a-zA-Za-яА-Я0-9\"\\s-])*)"
	// остаток
	regexpBalik := "(Баланс:)"

	var regp []string
	regp = append(regp, regexpCard)
	regp = append(regp, regexpDate)
	regp = append(regp, regexpTime)
	regp = append(regp, regexpKeyWord)
	regp = append(regp, regexpSum)
	regp = append(regp, regexpW)
	regp = append(regp, regexpBalik)
	regp = append(regp, regexpSum) // остаток

	var text string
	for i, reg := range regp {
		if i != len(regp)-1 {
			text += reg + regexpSeparator
		} else {
			text += reg
		}
	}

	re := regexp.MustCompile(text)

	fmt.Println(text)

	// выделение ключевых слоов
	if re.MatchString(message) {
		/*
			В массиве rez содержатся данные:
			0 - полностью сообщение
			1 - тип карты с указанием 4-х последних цифр карты
			2 - тип карты
			3 - дата совершения операции в формате dd.mm.yy
			4 - день совершения операции
			5 - месяц совершения операции
			6 - 2-е последние цифры года совершения операции
			7 - время совершения операции в формате hh:mm
			8 - час совершения операции
			9 - минуты совершения операции
			10 - операция
			11 - операция
			12 - сумма операции с указанием валюты
			13 - целая часть суммы
			14 - точка + дробная часть суммы
			15 - магазин/?? где совершили покупку
			16 - последняя буква магзина где совершили покупку
			17 - cлово "Баланс:"
			18 - остаток на карте с указанием валюты
			19 - целая часть остатка на карте
			20 - точка + дробная часть остатка на карте
		*/
		rez := re.FindStringSubmatch(message)

		data.keyWord = rez[10]
		data.amount = rez[13]
		data.targetName = rez[15]
		data.date = rez[3]
	}
	return data

}

// AddCosts добавляет данные из data в страницу учета расходов
func AddCosts(data jsonStruct) {

	// выбор столбца
	// выделение месяца
	mas_date := strings.Split(data.date, ".")
	// выделение месяца
	month := mas_date[1]

	var col_cost string
	// выбор столбца для записи
	switch month {
	case "01":
		col_cost = "B"
	case "02":
		col_cost = "C"
	case "03":
		col_cost = "D"
	case "04":
		col_cost = "E"
	case "05":
		col_cost = "F"
	case "06":
		col_cost = "G"
	case "07":
		col_cost = "H"
	case "08":
		col_cost = "I"
	case "09":
		col_cost = "J"
	case "10":
		col_cost = "K"
	case "11":
		col_cost = "L"
	case "12":
		col_cost = "M"
	}

	//выбор строки
	var row_cost int
	switch data.keyWord {
	case "покупка":
		row_cost = 2
	case "списание":
		row_cost = 3
	case "оплата":
		row_cost = 4
	case "перевод":
		row_cost = 5
	}

	//выделение ячейки
	range_cost := "Costs!" + col_cost + string(row_cost)
	fmt.Println("Rezultat costs = " + range_cost)

	//считать значение выделенной ячейки
	// resp, err := srvSpreadsheets.Values.Get(spreadsheetId, range_cost).Context(ctx).Do()
	// summ_cost = resp + summ
}

// AddAim добавляет данные из data в страницу учета доходов
func AddAim(data jsonStruct) {

	// выбор столбца
	// выделение месяца
	mas_date := strings.Split(data.date, ".")
	month := mas_date[1]

	var col_aim string
	switch month {
	case "01":
		col_aim = "B"
	case "02":
		col_aim = "C"
	case "03":
		col_aim = "D"
	case "04":
		col_aim = "E"
	case "05":
		col_aim = "F"
	case "06":
		col_aim = "G"
	case "07":
		col_aim = "H"
	case "08":
		col_aim = "I"
	case "09":
		col_aim = "J"
	case "10":
		col_aim = "K"
	case "11":
		col_aim = "L"
	case "12":
		col_aim = "M"
	}

	//выбор строки
	// определение типа дохода
	var row_aim int
	switch data.keyWord {
	case "зарплата":
		row_aim = 2
	case "пенсия":
		row_aim = 3
	case "зачисление":
		row_aim = 4
	case "аванс":
		row_aim = 5
	}

	//выделение ячейки
	range_aim := fmt.Sprintf("Доходы!%s%d", col_aim, row_aim)

	majorDimension := "COLUMNS"

	//считать значение выделенной ячейки
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, range_aim).MajorDimension(majorDimension).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	var sum, summ float64
	str := fmt.Sprintf("%v", resp.Values[0])
	sum, _ = strconv.ParseFloat(str, 64)
	summ, _ = strconv.ParseFloat(data.amount, 64)
	summ_aim := sum + summ

	d, err := resp.MarshalJSON()
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}
	fmt.Println(d[70], d[71], d[72])

	var f [][]interface{}
	f = summ_aim
	rb := &sheets.ValueRange{Values: f}

	/*	resp2, err := srv.Spreadsheets.Values.Update(spreadsheetId, range_aim, rb).Context(ctx).Do()
		if err != nil {
			log.Fatal(err)
		}
	*/
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

		//fmt.Printf("%s\n", update.Message.From.ID)

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
		if update.Message.Command() != "" {
			switch update.Message.Command() {
			case "auth": // авторизация пользователя
				ChatID = update.Message.Chat.ID
				UzverID = update.Message.From.ID
				go initGoogle()
				isReadAuth = true
				break
			case "file":
			}

		} else {
			// считывание сообщений
			sms := update.Message.Text

			fmt.Println("Было введено: " + sms)

			// парсинг ключевых слов
			data := ParseMess(sms)

			fmt.Print("Было поллучено из сообщения: ")
			fmt.Println(data)

			// вызов определенных функций
			if data.keyWord == " " {
				SendMessageToTelegram("Ошибка в сообщении")
			} else {
				switch data.keyWord {
				case "зачисление", "зачисление зарплаты", "пополнение":
					AddAim(data)
				case "списание", "снятие":
					AddCosts(data)
				}
			}
		} // else

	} // for update := range updates
}
