package main

// Config структура файла с настройками Telegram
type Config struct {
	TelegramBotToken string
}

//структура json
type jsonStruct struct {
	keyWord    string //ключевое слово из сообщения
	amount     string //сумма, указанная в сообщении
	targetName string //имя цели, если ключевое слово-цель
	date       string //дата, указанная в сообщении
}

type grebanniyJson struct {
	value []Value `json:"values"`
}

type Value struct {
	str []string
}
