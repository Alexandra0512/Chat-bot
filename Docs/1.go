package main

import (
	"fmt"
	"strings"
)

func main() {

	var (
		key     string
		sms     string
		mas_sms []string
		mas_key []string
	)

	key = "Справка spravka Зачисление zachislenie Пополнение popolnenie Поступление postuplenie Списание Spisanie Покупка Perevod Оплата	Oplata Цель tsel Баланс	balans Остаток	Ostatok Выписка Vypiska Редектировать Redektirovat Изменить Izmenit Кошелек	Koshelek Долг Dolg Создать Sozdat Выход vykhod"
	sms = "Schet *1096: postuplenie zarabotnoy plati 5605.25 RUB; 02.03.2018 05:36:31; Dostupno 5938.37 RUB. Detali v mobilnom banke: vtb.ru/app"

	mas_sms = strings.Split(sms, " ") //переход от строки к массиву
	//	fmt.Printf("%q\n", mas_sms)
	mas_key = strings.Split(key, " ")
	//	fmt.Printf("%q\n", mas_key)

	for i := 0; i < len(mas_sms); i++ {
		for j := 0; j < len(mas_key); j++ {
			if mas_sms[i] == mas_key[j] {
				fmt.Printf("%q\n", mas_sms[i])
			}
		}
	}
}
