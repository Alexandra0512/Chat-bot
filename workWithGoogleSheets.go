package main

import (
	"fmt"
	"log"
	"strconv"
)

// Результат труда Фёдоровой Александры
// getKeywordID получение ID ключевого слова из таблицы.
// На вход поступает ключевое слова
// Возвращает id ключевого слова или 0 - если ID не найдено
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

// Результат труда Фёдоровой Александры
// getAnswerID Возвращает ID ответа по ID ключевого слова*/
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

// Результат труда Фёдоровой Александры
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
