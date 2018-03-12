package main

import (
    "fmt"
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
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	sheets "google.golang.org/api/sheets/v4"
	"gopkg.in/telegram-bot-api.v4"
)


func main() {
var (						
	key string				
	sms string	
    mas_sms	[] string
	mas_key	[] string
	)
	for update := range updates {
		// Создав структуру - можно её отправить обратно боту

		// Название таблицы, куда будут заносится данные
		// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit
		spreadsheetId := "1XbhJ785LzQ2O713foKvOchjzzg8M-rPtePOfAOvU83Y"	
		// указывается страница и диапозон
		// сейчас происходит чтение таблицы с указанием пользователей
		readRange := "TableKeyWords!A2:36"
		resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
		if err != nil {
			log.Fatalf("Unable to retrieve data from sheet. %v", err)
		}

		key := ""
		if len(resp.Values) > 0 {
			// считывание ключевых слов из таблицы TableKeyWords
			for _, row := range resp.Values {
				// Print columns A and E, which correspond to indices 0 and 4.
				fmt.Printf("%s -> %s\n", row[0], row[1])
				s := []string{row[0].(string), row[1].(string), "\n"}
				key += strings.Join(s, ",")
			}
		} else {
			fmt.Print("No data found.")
		}
//    key= "Справка spravka Зачисление zachislenie Пополнение popolnenie Поступление postuplenie Списание Spisanie Покупка Perevod Оплата	Oplata Цель tsel Баланс	balans Остаток	Ostatok Выписка Vypiska Редектировать Redektirovat Изменить Izmenit Кошелек	Koshelek Долг Dolg Создать Sozdat Выход vykhod"
	sms= "Schet *1096: postuplenie zarabotnoy plati 5605.25 RUB; 02.03.2018 05:36:31; Dostupno 5938.37 RUB. Detali v mobilnom banke: vtb.ru/app"

	mas_sms= strings.Split(sms, " ") //переход от строки к массиву 
//	fmt.Printf("%q\n", mas_sms)
	mas_key= strings.Split(key, " ")
//	fmt.Printf("%q\n", mas_key)

		for i := 0; i < len(mas_sms); i++ {
			for j:=0; j < len(mas_key); j++ {
				if (mas_sms[i]==mas_key[j]){
				fmt.Printf("%q\n", mas_sms[i])
			}
			}
		}		
	
	}
}