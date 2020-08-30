/*
 * Copyright (c) 2020 GeoSonic. All rights reserved.
 */

package bot

import (
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/SevereCloud/vksdk/longpoll-user/v3"

	"github.com/SevereCloud/vksdk/api"
	"github.com/SevereCloud/vksdk/api/errors"
	"github.com/SevereCloud/vksdk/longpoll-user"
)

const Version = "1.2"

// Запускает аккаунт
func start(token string, triggerWord interface{}) {
	vk := api.NewVK(token)
	// Чтобы избежать конфликтов в отправке
	// запросов, установим лимит
	vk.Limit = api.LimitUserToken
	vk.UserAgent += ", toDelete/" + Version + " (+https://github.com/geosonic/todelete)"
	lp, err := longpoll.NewLongpoll(vk, longpoll.ReceiveAttachments)
	if err != nil {
		if errors.GetType(err) == errors.Auth {
			log.Println("Account run failed: invalid token")
		} else {
			log.Printf("Account run failed [*%s]: %v\n", token[len(token)-4:], err.Error())
		}
		return
	}

	regexp1, _ := parse(triggerWord)

	if regexp1 == nil {
		log.Printf("Account run failed [*%s]: Must be trigger word string or []string.\n", token[len(token)-4:])
	}

	w := wrapper.NewWrapper(lp)

	acc := token[len(token)-4:]

	w.OnNewMessage(func(message wrapper.NewMessage) {
		/* TODO: Лог необходимо переработать */
		// Проверяем только свои сообщения
		if !message.Flags.Has(wrapper.Outbox) {
			return
		}

		// Проверяем сообщение
		result := regexp1.FindStringSubmatch(message.Text)

		if result == nil {
			return
		}

		var toDelete struct {
			replace bool
			count   int
		}

		// Если сообщение закончится на триггер + "-",
		// значит сообщение будет отредактировано, также
		// после "-" можно написать кол-во сообщений
		// для удаления
		if result[1] == "-" {
			toDelete.replace = true
		}

		toDelete.count, err = strconv.Atoi(result[2])

		if err != nil {
			toDelete.count = 1
		}

		if toDelete.replace {
			log.Printf("Delete replace in *%s (%d)\n", acc, toDelete.count)
			// Получение сообщений с помощью execute
			messages, err := GetMessages(vk, toDelete.count+1, message.PeerID)
			if err != nil {
				log.Printf("[*%s] Error getting messages (%v)", acc, err.Error())
				return
			}

			if len(messages) == 0 {
				log.Printf("[*%s] Not found messages for edit.", acc)
				return
			}

			// Сортировка списка (в нашем случае он будет перевёрнут)
			sort.Ints(messages)

			var count int

			for _, v := range messages {
				// Проверка на сообщение, вызвавшего
				// удаление сообщений, лично я не вижу
				// смысла редактировать такое сообщение
				if v != message.MessageID {
					_, err := vk.MessagesEdit(api.Params{"peer_id": message.PeerID, "message_id": v, "message": "ᅠ"})
					if err == nil {
						count++
						log.Printf("Edited %v messages\n", count)
					} else if errors.GetType(err) == errors.Captcha {
						log.Println(err)
						break
					}

					// Задержка для корректного удаления
					if len(messages) > 2 {
						time.Sleep(time.Second / 2)
					}
				}
			}

			log.Printf("[*%s] %d of %d messages edited.\n", acc, count, len(messages))

			for i := 0; i < 10; i++ {
				_, err = vk.MessagesDelete(api.Params{"message_ids": messages, "delete_for_all": 1})
				if err == nil {
					log.Printf("[*%s] %d messages deleted!\n", acc, len(messages))
					break
				}
				log.Printf("[*%s] Error deleting, trying (%v)\n", acc, i)
			}
		} else {
			log.Printf("[*%s] Delete %d messages\n", acc, toDelete.count)
			// Удаление сообщений с помощью execute
			deleted, err := DeleteExec(vk, toDelete.count+1, message.PeerID)
			if err != nil {
				log.Printf("[*%s] Error deleting messages! (%v)\n", acc, err.Error())
			}
			log.Printf("[*%s] Deleted %d messages\n", acc, deleted)
		}
	})

	// Запуск аккаунта
	for {
		log.Printf("Run *%s\n", acc)
		err := lp.Run()
		if err != nil {
			log.Printf("[*%s] LongPoll connecting error: %v, trying...", acc, err)
		}
		time.Sleep(time.Second * 10)
	}
}

// Запускает аккаунты параллельно
func StartAccounts(accounts map[string]interface{}) {
	for k, v := range accounts {
		if k != "" && v != "" {
			go start(k, v)
		}
	}
	select {}
}
