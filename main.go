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
	var accounts map[string]bot.Config

	// Путь к конфигу
	configPath := flag.String("configPath", "config.json", "path to config")
	flag.Parse()

L1:
	file, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Fatalln("Error reading file:", err)
	}

	if !json.Valid(file) {
		log.Fatalln("Incorrect config file.")
	}

	err = json.Unmarshal(file, &accounts)
	if err != nil {
		// Если конфиг будет битым - он запаникует, так что норм
		bot.TranspileConfig(*configPath)
		goto L1 // Знаю, что goto это плохо, но так красиво
	}

	if len(accounts) == 0 {
		log.Fatalln("Accounts not found!")
	}

	// Запуск аккаунтов
	bot.StartAccounts(accounts)
}
