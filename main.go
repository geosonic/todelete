/*
 * Copyright (c) 2020 GeoSonic. All rights reserved.
 */

package main

import (
	"deleter/bot"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
)

func main() {
	/* Аккаунты теперь должны находиться в config.json */
	fmt.Printf("To Delete %s by GeoSonic for %s_%s\n", bot.Version, runtime.GOOS, runtime.GOARCH)
	var accounts map[string]interface{}

	// Путь к конфигу
	configPath := flag.String("configPath", "config.json", "path to config")
	flag.Parse()

	file, err := ioutil.ReadFile(*configPath)
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

	// Запуск аккаунтов
	bot.StartAccounts(accounts)
}
