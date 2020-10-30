/*
 * Copyright (c) 2020 GeoSonic. All rights reserved.
 */

package bot

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/SevereCloud/vksdk/v2/longpoll-user/v3"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/longpoll-user"
)

const Version = "2.0.0"

// Запускает аккаунт
func start(token string, config Config, vk *api.VK) {
	// Вообще я задумывался о возможности выключить редактирование,
	// но я пока решил пока этого не делать
	if config.Separator == "" {
		config.Separator = "-"
	}

	lp, err := longpoll.NewLongPoll(vk, longpoll.ReceiveAttachments)
	lp.Goroutine(false)
	if err != nil {
		if errors.Is(err, api.ErrAuth) {
			log.Println("Account run failed: invalid token")
		} else {
			log.Printf("Account run failed [*%s]: %v\n", token[len(token)-4:], err.Error())
		}
		return
	}

	// language=regexp
	regexp1 := regexp.MustCompile(fmt.Sprintf("(?i)^(?:%s)(%s+)?([0-9]+)?", strings.Join(config.Keywords, "|"), regexp.QuoteMeta(config.Separator)))

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

		// Если сообщение закончится на config.Separator,
		// значит сообщение будет отредактировано, также
		// после "-" можно написать кол-во сообщений
		// для удаления
		if strings.Contains(result[1], config.Separator) {
			toDelete.replace = true
		}

		// Пусть будет ключевое слово "хм" и маркер редактирования "+"
		// Перед удалением мы написали сообщения "Привет", "Как дела?", "Что делаешь?" по порядку
		//
		// Если напишем "хм", то удалится "Что делаешь?"
		// Если напишем "хм2", то удалится "Как дела?" и "Что делаешь?"
		// Если напишем "хм+, то сообщение "Что делаешь" отредактируется, и удалится
		// Если напишем "хм+2", то "Как дела?" и "Что делаешь?" отредактируются и удалятся
		//
		// Маркеры как количество для удаления:
		// Если напишем "хм+, то сообщение "Что делаешь" отредактируется, и удалится
		// Если напишем "хм++", то "Как дела?" и "Что делаешь?" отредактируются и удалятся
		// Если напишем "хм+++", то удалятся и отредактируются все эти 3 сообщения
		//
		// И т.д.
		//
		// По умолчанию удаляется запрашиваемое количество + 1 сообщений, т.к. нужно удалить ещё сообщение с триггером
		// Кстати, если написать, например, "хм0", то он удалит только сообщение, которое вызвало удаление

		sepCount := strings.Count(result[1], config.Separator)

		toDelete.count, err = strconv.Atoi(result[2])

		if err != nil {
			if toDelete.replace && config.CountSeparator && sepCount > 1 {
				toDelete.count = sepCount
			} else {
				toDelete.count = 1
			}
		}

		if config.DeleteTriggerFast {
			if config.EditTrigger {
				_, err := vk.MessagesEdit(api.Params{"peer_id": message.PeerID, "message_id": message.MessageID, "message": config.ToEditString})
				log.Println("Error editing trigger message: ", err)
			}
			vk.MessagesDelete(api.Params{"message_ids": message.MessageID, "delete_for_all": 1})
		}

		if toDelete.replace {
			log.Printf("Delete replace in *%s (%d)\n", acc, toDelete.count)
			// Получение сообщений с помощью execute
			messages, err := GetMessages(vk, toDelete.count, message.PeerID)
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
				if v == message.MessageID && !(config.EditTrigger && config.DeleteTriggerFast) {
					continue
				}

				_, err := vk.MessagesEdit(api.Params{"peer_id": message.PeerID, "message_id": v, "message": config.ToEditString})
				if errors.Is(err, api.ErrCaptcha) {
					log.Println(err)
					break
				}

				count++
				log.Printf("Edited %v messages\n", count)

				// Задержка для корректного удаления
				if len(messages) > 2 {
					time.Sleep(time.Second / 2)
				}
			}

			log.Printf("[*%s] %d of %d messages edited.\n", acc, count, len(messages))

			for i := 0; i < 10; i++ {
				_, err = vk.MessagesDelete(api.Params{"message_ids": messages, "delete_for_all": 1})
				if err == nil {
					log.Printf("[*%s] %d messages deleted!\n", acc, len(messages))
					break
				}
				fmt.Println(err)
				log.Printf("[*%s] Error deleting, trying (%v)\n", acc, i)
			}
			return
		}

		log.Printf("[*%s] Delete %d messages\n", acc, toDelete.count)
		// Удаление сообщений с помощью execute
		deleted, err := DeleteExec(vk, toDelete.count, message.PeerID)
		if err != nil {
			log.Printf("[*%s] Error deleting messages! (%v)\n", acc, err.Error())
		}
		log.Printf("[*%s] Deleted %d messages\n", acc, deleted)
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
func StartAccounts(accounts map[string]Config) {
	for k, v := range accounts {
		if k != "" {
			vk := api.NewVK(k)
			// Чтобы избежать конфликтов в отправке
			// запросов, установим лимит
			vk.Limit = api.LimitUserToken
			vk.UserAgent += ", toDelete/" + Version + " (+https://github.com/geosonic/todelete)"

			go start(k, v, vk)
		}
	}
	select {}
}
