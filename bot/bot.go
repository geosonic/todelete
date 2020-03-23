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
		log.Fatal(fmt.Sprintf("Account *%v: %v", token[len(token)-4:], err))
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
			fmt.Printf("Delete replace in *%v (%v)\n", token[len(token)-4:], count)
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
			for i := 0; i < 10; i++ {
				_, err = vk.MessagesDelete(api.Params{"message_ids": messages, "delete_for_all": 1})
				if err == nil {
					break
				}
			}

		} else {
			fmt.Printf("Delete in *%v (%v)\n", token[len(token)-4:], count)
			// Удаление сообщений с помощью execute
			DeleteExec(vk, count+1, message.PeerID)
		}

		return
	})

	// Запуск и автоподнятие
	for {
		fmt.Printf("Run *%v", token[len(token)-4:])
		_ = lp.Run()
		time.Sleep(time.Second * 10)
	}
}
