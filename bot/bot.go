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
		log.Fatal(err)
	}

	regexp1 := regexp.MustCompile(fmt.Sprintf("^%v(-)?([0-9]+)?", strings.ToLower(triggerWord)))

	w := wrapper.NewWrapper(&lp)

	w.OnNewMessage(func(message wrapper.NewMessage) {
		// Проверяем только свои сообщения
		if !message.Flags.Has(wrapper.Outbox) {
			return
		}

		message.Text = strings.ToLower(message.Text)

		// Проверяем сообщение
		result := regexp1.FindStringSubmatch(message.Text)

		var (
			toDeleteReplace bool
			count           int
		)

		if result == nil {
			return
		}

		if result[1] == "-" {
			toDeleteReplace = true
		}
		count, err = strconv.Atoi(result[2])
		if err != nil {
			count = 1
		}

		if toDeleteReplace {
			// Получение сообщений с помощью execute
			messages := GetMessages(vk, count+1, message.PeerID)

			// Переворачиваем список
			for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
				messages[i], messages[j] = messages[j], messages[i]
			}

			for _, v := range messages {
				if v != message.MessageID {
					_, err := vk.MessagesEdit(api.Params{"peer_id": message.PeerID, "message_id": v, "message": "ᅠ"})

					if errors.GetType(err) == errors.Captcha {
						break
					}

					// Задержка для корректного удаления
					time.Sleep(time.Millisecond * 500)
				}
			}
			messages = append(messages, message.MessageID)

			for i := 0; i < 10; i++ {
				_, err = vk.MessagesDelete(api.Params{"message_ids": ToArray(messages), "delete_for_all": 1})
				if err == nil {
					break
				}
			}

		} else {
			// Удаление сообщений с помощью execute
			DeleteExec(vk, count+1, message.PeerID)
		}

		return
	})

	// Запуск и автоподнятие
	for {
		_ = lp.Run()
		time.Sleep(time.Second * 10)
	}
}
