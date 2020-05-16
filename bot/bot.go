/*
 * Copyright (c) 2020 GeoSonic. All rights reserved.
 */

package bot

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/SevereCloud/vksdk/longpoll-user/v3"

	"github.com/SevereCloud/vksdk/api"
	"github.com/SevereCloud/vksdk/api/errors"
	"github.com/SevereCloud/vksdk/longpoll-user"
)

// Запускает аккаунт
func start(token, triggerWord string) {
	vk := api.NewVK(token)
	vk.Limit = api.LimitUserToken
	lp, err := longpoll.NewLongpoll(vk, 2)
	if err != nil {
		log.Printf("Account *%v: %v\n", token[len(token)-4:], err)
		return
	}

	regexp1 := regexp.MustCompile(fmt.Sprintf("^%v(-)?([0-9]+)?", strings.ToLower(triggerWord)))

	w := wrapper.NewWrapper(lp)

	acc := token[len(token)-4:]

	w.OnNewMessage(func(message wrapper.NewMessage) {
		/* TODO: Лог необходимо переработать */
		// Проверяем только свои сообщения
		if !message.Flags.Has(wrapper.Outbox) {
			return
		}

		message.Text = strings.ToLower(message.Text)

		// Проверяем сообщение
		result := regexp1.FindStringSubmatch(message.Text)

		if result == nil {
			return
		}

		var info struct {
			replace bool
			count   int
		}

		if result[1] == "-" {
			info.replace = true
		}

		info.count, err = strconv.Atoi(result[2])

		if err != nil {
			info.count = 1
		}

		if info.replace {
			log.Printf("Delete replace in *%v (%v)\n", acc, info.count)
			// Получение сообщений с помощью execute
			messages, err := GetMessages(vk, info.count+1, message.PeerID)
			if err != nil {
				log.Printf("[*%v] Error getting messages (%v)", acc, err.Error())
			}

			// Переворачиваем список, если вы хотите, чтобы сообщения
			// удалялись в другом порядке, этот цикл нужно убрать
			for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
				messages[i], messages[j] = messages[j], messages[i]
			}

			var count int

			for _, v := range messages {
				if v != message.MessageID {
					_, err := vk.MessagesEdit(api.Params{"peer_id": message.PeerID, "message_id": v, "message": "ᅠ"})
					if err == nil {
						count++
						log.Printf("Edited %v messages\n", count)
					} else {
						if errors.GetType(err) == errors.Captcha {
							log.Println(err)
							break
						}
					}

					// Задержка для корректного удаления
					if len(messages) > 2 {
						time.Sleep(time.Second / 2)
					}
				}
			}

			log.Printf("[*%v] %v of %v messages edited.\n", acc, count, len(messages))

			for i := 0; i < 10; i++ {
				_, err = vk.MessagesDelete(api.Params{"message_ids": messages, "delete_for_all": 1})
				if err == nil {
					log.Printf("[*%v] %v messages deleted!\n", acc, len(messages))
					break
				}
				log.Printf("[*%v] Error deleting, trying (%v)\n", acc, i)
			}

		} else {
			log.Printf("[*%v] Delete %v messages\n", acc, info.count)
			// Удаление сообщений с помощью execute
			err = DeleteExec(vk, info.count+1, message.PeerID)
			if err != nil {
				log.Printf("[*%v] Error deleting messages! (%v)\n", acc, err.Error())
			}
		}
		return
	})

	// Запуск и автоподнятие
	for {
		log.Printf("Run *%v\n", acc)
		_ = lp.Run()
		time.Sleep(time.Second * 10)
	}
}

// Запускает аккаунты параллельно
func StartAccounts(accounts map[string]string) {
	for k, v := range accounts {
		if k != "" && v != "" {
			go start(k, v)
		}
	}
	select {}
}
