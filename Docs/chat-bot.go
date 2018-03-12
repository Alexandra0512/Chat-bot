package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strconv"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	sheets "google.golang.org/api/sheets/v4"
	"gopkg.in/telegram-bot-api.v4"
)

type Config struct {
	TelegramBotToken string
}

var srv *sheets.Service
var spreadsheetId string = "1XbhJ785LzQ2O713foKvOchjzzg8M-rPtePOfAOvU83Y"

// getClient uses a Context and Config to retrieve a Token
// then generate a getCliClient. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok := getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("sheets.googleapis.com-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

//Саша
/*Возвращает ID ключевого слова из таблицы
0-если ID не найдено*/
func getKeywordID(keyword string) (keywordID int) {
	//обращаемся к таблице TableKeyWords
	readRange := "TableKeyWords!A2:19"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet. %v", err)
	}
	//ищем id ключевого слова
	if len(resp.Values) > 0 {
		for _, row := range resp.Values {
			if keyword == row[1] {
				//fmt.Print("Yes: ", row[0])
				num, err := strconv.Atoi(row[0].(string))
				if err == nil {
					return num
				} else {
					fmt.Println("Число имеет неверный формат")
					return 0
				}

			}
		}
	} else {
		fmt.Print("No data found.")

	}
	fmt.Print("No")
	return 0
}

/*Возвращает ID ответа по ID ключевого слова*/
func getAnswerID(keywordID int) (answerID int) {
	//обращаемся к таблице TableKeyAnsw
	readRange := "TableKeyAnsw!A2:19"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet. %v", err)
	}
	//ищем id ответа
	if len(resp.Values) > 0 {
		for _, row := range resp.Values {
			num, err := strconv.Atoi(row[1].(string))
			if err != nil {
				fmt.Println("Число в таблице имеет неверный формат")
				return 0
			}
			if keywordID == num {
				//fmt.Print("Yes: ", row[0])
				num, err := strconv.Atoi(row[0].(string))
				if err == nil {
					return num
				} else {
					fmt.Println("Число имеет неверный формат")
					return 0
				}

			}
		}
	} else {
		fmt.Print("No data found.")

	}
	fmt.Print("No")
	return 5
}

/*Ищет ответ по Id ответа*/
func getAnswerByID(answerID int) (answer string) {
	//обращаемся к таблице TableAnswer
	readRange := "TableAnswer!A2:19"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet. %v", err)
	}
	//ищем id ответа
	if len(resp.Values) > 0 {
		for _, row := range resp.Values {
			num, err := strconv.Atoi(row[0].(string))
			if err != nil {
				fmt.Println("Число в таблице имеет неверный формат")
				return
			}
			if answerID == num {
				return row[1].(string)
			}
		}
	} else {
		fmt.Print("No data found.")

	}
	fmt.Print("No")
	return "Не найдено ответа по данному ID"
}

/*Отправляет сообщение в телеграмме по ID сообщения*/
func sendMessage(answer string) {
	//bot.Send(answer)
}

func main() {
	ctx := context.Background()
	// Получение доступа к Google Sheets

	b, err := ioutil.ReadFile("config_google.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/sheets.googleapis.com-go-quickstart.json
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := getClient(ctx, config)

	srv, err = sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets Client %v", err)
	}

	// Получение доступа к Telegram
	file, _ := os.Open("config_telegram.json")
	decoder := json.NewDecoder(file)
	configuration := Config{}
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(configuration.TelegramBotToken)

	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Panic(err)
	}

	// В канал updates будут приходить все новые сообщения.
	for update := range updates {
		// Создав структуру - можно её отправить обратно боту

		//Саша
		message := update.Message.Text //считываем сообщение
		fmt.Printf("%s\n", message)
		//здесь должны выделить ключевое слово
		//Аня Ч.

		//Саша
		keywordID := getKeywordID(message) //получаем id ключевого слова
		fmt.Printf("keywordID is %d\n", keywordID)
		fmt.Printf("answerID is %d\n", getAnswerID(keywordID))
		fmt.Printf("answer is %s\n", getAnswerByID(getAnswerID(keywordID)))
		// формирование сообщения
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, getAnswerByID(getAnswerID(keywordID)))

		// вывод сообщения в телеграмме
		bot.Send(msg)
		//keywordID = getKeywordID(msg)
		// Название таблицы, куда будут заносится данные
		// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit

		// указывается страница и диапозон
		// сейчас происходит чтение таблицы с указанием пользователей

		//	text := ""
		/*	if len(resp.Values) > 0 {
			// считывание пользователей из таблицы TableUsers
			for _, row := range resp.Values {
				// Print coError 403: Project unknown (#9100848745) has been deleted., forbiddenlumns A and E, which correspond to indices 0 and 4.
				fmt.Printf("%s -> %s\n", row[0], row[1])
				s := []string{row[0].(string), row[1].(string), "\n"}
				text += strings.Join(s, ",")
			}*/

		// формирование сообщения
		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)

		// вывод сообщения в телеграмме
		//bot.Send(keywordID)
	}
}
