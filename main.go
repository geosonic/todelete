/*
 * Copyleft (ↄ) 2020, Geosonic
 */

package main

import (
	"deleter/bot"
	"encoding/json"
	"io/ioutil"
	"log"
)

func main() {
	/* Аккаунты теперь должны находиться в config.json */
	var accounts map[string]string
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(file, &accounts)
	if err != nil {
		panic(err)
	}

	if len(accounts) == 0 {
		log.Fatalln("Accounts not found!")
	}

	// Функция запуска аккаунтов
	bot.StartAccounts(accounts)
}
