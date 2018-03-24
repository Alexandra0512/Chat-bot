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

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	drive "google.golang.org/api/drive/v3"
	sheets "google.golang.org/api/sheets/v4"
)

var (
	AuthCode      chan string = make(chan string)
	spreadsheetId string
	srv           *sheets.Service
	ctx           context.Context
)

// initGoogle Инициализация связи с Google Drive и Google Sheets
func initGoogle() {
	ctx = context.Background()

	// чтение файлов настроек от гугл
	b, err := ioutil.ReadFile("config_google.json")
	if err != nil {
		log.Fatalf("Невозможно прочитать файл настроек google: %v", err)
		SendMessageToTelegram("Невозможно прочитать файл настроек google, обратитесь администратору")
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/sheets.googleapis.com-go-quickstart.json
	config, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Невозмоожно получить данный из файла настроек google: %v", err)
	}

	// создание клиента
	client := getClient(ctx, config)

	// получение доступа к Google Drive
	srvDrive, err = drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve drive Client %v", err)
	}

	// получение доступа к Google Sheets
	srv, err = sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets Client %v", err)
	}

	// просмотр файлов на гугл диске
	// выделение у файла id, название (name) и тип файла (mimeType)
	r, er := srvDrive.Files.List().PageSize(10).Fields("nextPageToken, files(id, name, mimeType)").Do()
	if er != nil {
		log.Fatalf("Не удаётся получить файлы с Google Drive: %v", err)
	}

	Table(r, ctx)
}

// Table проверка есть ли на гугл диске пользователя таблица с учетом бюджета.
// Если таблицы нет, то она создаётся
// r - список файлов на гугл диске пользователя
// ctx - ?????
func Table(r *drive.FileList, ctx context.Context) {
	// такое расширение у таблиц в гуггл диске
	typeSpreadsheet := "application/vnd.google-apps.spreadsheet"
	link := "https://docs.google.com/spreadsheets/d/"

	countFile := 0
	isTableNoExist := false

	if len(r.Files) > 0 {
		// поиск среди файлов пользователя электронной таблицы с именем "Бюджет"
		for _, i := range r.Files {
			if i.Name == "Бюджет" && i.MimeType == typeSpreadsheet {
				spreadsheetId = i.Id
				break
			} else {
				countFile++
			}
		} // for
	} else {
		isTableNoExist = true
	}

	// Если среди файлов таблица не найдена, значит надо её создать
	if countFile == len(r.Files) {
		isTableNoExist = true
	}

	// создание таблицы, если она отсутствует на диске пользователя
	if isTableNoExist {
		table := StructTableSheets()
		resp, err := srv.Spreadsheets.Create(table).Context(ctx).Do()
		if err != nil {
			log.Fatal(err)
		}
		spreadsheetId = resp.SpreadsheetId
	}

	SendMessageToTelegram("Ссылка на таблицу:\n" + link + spreadsheetId)
	isReadAuth = false

}

// StructTableSheets создание структуры листов таблицы
// Возвращает структуру листов таблицы (названия листов, их порядок)
func StructTableSheets() *sheets.Spreadsheet {
	file, _ := os.Open("table.json")
	defer file.Close()

	jsonDataTable, _ := ioutil.ReadAll(file)
	var table *sheets.Spreadsheet
	json.Unmarshal(jsonDataTable, &table)

	return table
}

// getClient использование Context и Config для получения Token
// на основе которго генерируется Client. Возвращает ссылку на Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {

	// Создание файла учетных данных
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Не удаётся получить доступ к файлу учётных записей. %v", err)
	}

	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb использование Config для запроса Token.
// Возвращает Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	text := "Откройте ссылку для авторизации в Google: \n" + authURL
	SendMessageToTelegram(text)

	tok, err := config.Exchange(oauth2.NoContext, <-AuthCode)
	if err != nil {
		SendMessageToTelegram("Не удаётся получить Token  ")
	}

	return tok

}

// tokenCacheFile создание пути к файлу учетных данных (имени файла).
// Возвращает сгенерированный путь к файлу учетных данных (имени файла).
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("quickstart"+string(UzverID)+".json")), err
}

// tokenFromFile извлечение Token из пути к файлу.
// Возвращает Token и ошибку, обнаруженную во время чтения файла.
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

// saveToken создание файла и сохранение в нём Token.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Файл учетных данных сохранён в : %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Не удаётся кэшировать oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
