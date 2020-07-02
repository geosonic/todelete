/*
 * Copyright (c) 2020 GeoSonic. All rights reserved.
 */

package main

import (
	"deleter/bot"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
)

func main() {
	/* Аккаунты теперь должны находиться в config.json */
	fmt.Printf("To Delete %v by GeoSonic for %v_%v\n", bot.Version, runtime.GOOS, runtime.GOARCH)
	var accounts map[string]interface{}
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalln("Error reading file:", err)
	}

	if !json.Valid(file) {
		log.Fatalln("Incorrect config file.")
	}

	err = json.Unmarshal(file, &accounts)
	if err != nil {
		log.Fatalln("Can't unmarshal json:", err)
	}

	if len(accounts) == 0 {
		log.Fatalln("Accounts not found!")
	}

	// Функция запуска аккаунтов
	bot.StartAccounts(accounts)
}
