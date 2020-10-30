/*
 * Copyright (c) 2020 GeoSonic. All rights reserved.
 */

package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func parseKeyWord(keyword interface{}) []string {
	fmt.Println(keyword)
	switch tr := keyword.(type) {
	case string:
		return []string{tr}
	case []interface{}:
		var keyWords = make([]string, 0, len(tr))
		for _, v := range tr {
			switch k := v.(type) {
			case string:
				keyWords = append(keyWords, k)
			default:
				panic(fmt.Errorf("i want string or []string, no %T", keyword))
			}
		}

		return keyWords
	default:
		panic(fmt.Errorf("i want string or []string, no %T", keyword))
	}
}

type Config struct {
	Keywords          []string `json:"keywords"`
	Separator         string   `json:"separator"`
	ToEditString      string   `json:"to_edit_string"`
	CountSeparator    bool     `json:"count_separator"`
	DeleteTriggerFast bool     `json:"delete_trigger_fast"`
	EditTrigger       bool     `json:"edit_trigger"`
}

// Переводит формат старого конфига в новый формат
// Совмещать типы из старого формата и нового никак нельзя!
func TranspileConfig(configPath string) {
	/* Аккаунты теперь должны находиться в config.json */
	var accounts map[string]interface{}

	file, err := ioutil.ReadFile(configPath)
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

	var newConfig = make(map[string]Config)

	for token, v := range accounts {
		config := Config{
			Keywords:          parseKeyWord(v),
			Separator:         "-", // по умолчанию задан такой триггер для редактирования
			ToEditString:      "ᅠ",
			CountSeparator:    false,
			DeleteTriggerFast: false,
			EditTrigger:       false,
		}
		newConfig[token] = config
	}
	content, err := json.MarshalIndent(newConfig, "", "\t")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(configPath, content, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
