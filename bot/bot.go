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

func Start(token, triggerWord string) {
	vk := api.Init(token)
	lp, err := longpoll.Init(vk, 2)
	if err != nil {
		log.Fatal(fmt.Sprintf("Account *%v: %v\n", token[len(token)-4:], err))
	}

	regexp1 := regexp.MustCompile(fmt.Sprintf("^%v(-)?([0-9]+)?", strings.ToLower(triggerWord)))

	w := wrapper.NewWrapper(&lp)

	acc := token[len(token)-4:]

	w.OnNewMessage(func(message wrapper.NewMessage) {
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

		var (
			toDeleteReplace bool
			count           int
		)

		if result[1] == "-" {
			toDeleteReplace = true
		}

		count, err = strconv.Atoi(result[2])

		if err != nil {
			count = 1
		}

		if count > 10000 {
			count = 2147483647
		}
		// TODO: Сделать всё это в горутину.
		if toDeleteReplace {
			fmt.Printf("Delete replace in *%v (%v)\n", acc, count)
			// Получение сообщений с помощью execute
			messages := GetMessages(vk, count+1, message.PeerID)

			// Переворачиваем список
			for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
				messages[i], messages[j] = messages[j], messages[i]
			}

			var count int

			for _, v := range messages {
				if v != message.MessageID {
					_, err := vk.MessagesEdit(api.Params{"peer_id": message.PeerID, "message_id": v, "message": "ᅠ"})
					if err == nil {
						count++
					}

					if errors.GetType(err) == errors.Captcha {
						log.Println(err)
						break
					}

					// Задержка для корректного удаления
					if len(messages) > 2 {
						time.Sleep(time.Second / 2)
					}
				}
			}
			log.Printf("[*%v] %v of %v messages edited.\n", count, len(messages), acc)
			for i := 0; i < 10; i++ {
				_, err = vk.MessagesDelete(api.Params{"message_ids": messages, "delete_for_all": 1})
				if err == nil {
					log.Printf("[*%v] %v messages deleted!\n", acc, len(messages))
					break
				}
				log.Printf("[*%v] Error deleting, trying (%v)\n", acc, i)
			}

		} else {
			log.Printf("[*%v] Delete %v messages\n", acc, count)
			// Удаление сообщений с помощью execute
			err = DeleteExec(vk, count+1, message.PeerID)
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
