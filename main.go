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
	fmt.Printf("To Delete 1.1 by GeoSonic for %v_%v\n", runtime.GOOS, runtime.GOARCH)
	var accounts map[string]interface{}
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalln("Error reading file:", err)
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
