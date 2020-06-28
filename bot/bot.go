/*
 * Copyright (c) 2020 GeoSonic. All rights reserved.
 */

package bot

import (
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/SevereCloud/vksdk/longpoll-user/v3"

	"github.com/SevereCloud/vksdk/api"
	"github.com/SevereCloud/vksdk/api/errors"
	"github.com/SevereCloud/vksdk/longpoll-user"
)

// Запускает аккаунт
func start(token string, triggerWord interface{}) {
	vk := api.NewVK(token)
	// Чтобы избежать конфликтов в отправке
	// запросов, установим лимит
	vk.Limit = api.LimitUserToken
	lp, err := longpoll.NewLongpoll(vk, longpoll.ReceiveAttachments)
	if err != nil {
		if errors.GetType(err) == errors.Auth {
			log.Println("Account run failed: invalid token")
		} else {
			log.Printf("Account run failed [*%v]: %v\n", token[len(token)-4:], err.Error())
		}
		return
	}

	regexp1, _ := parse(triggerWord)

	if regexp1 == nil {
		log.Printf("Account run failed [*%v]: Must be trigger word string or []string.\n", token[len(token)-4:])
	}

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

		var toDelete struct {
			replace bool
			count   int
		}

		if result[1] == "-" {
			toDelete.replace = true
		}

		toDelete.count, err = strconv.Atoi(result[2])

		if err != nil {
			toDelete.count = 1
		}

		if toDelete.replace {
			log.Printf("Delete replace in *%v (%v)\n", acc, toDelete.count)
			// Получение сообщений с помощью execute
			messages, err := GetMessages(vk, toDelete.count+1, message.PeerID)
			if err != nil {
				log.Printf("[*%v] Error getting messages (%v)", acc, err.Error())
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
					} else {
						// Если же мы не смогли отредактировать сообщение
						// из-за капчи, немедленно срываем цикл и
						// удаляем все сообщения, т.к. больше пытаться
						// редактировать не имеет смысла
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
			log.Printf("[*%v] Delete %v messages\n", acc, toDelete.count)
			// Удаление сообщений с помощью execute
			err = DeleteExec(vk, toDelete.count+1, message.PeerID)
			if err != nil {
				log.Printf("[*%v] Error deleting messages! (%v)\n", acc, err.Error())
			}
		}
	})

	// Запуск и автоподнятие
	for {
		log.Printf("Run *%v\n", acc)
		_ = lp.Run()
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
