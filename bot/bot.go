package bot

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/SevereCloud/vksdk/api"
	"github.com/SevereCloud/vksdk/api/errors"
	"github.com/SevereCloud/vksdk/longpoll-user"
)

type Message struct {
	Event     int
	ID        int
	Flags     int
	PeerID    int
	Timestamp int
	Text      string
}

func Start(token, triggerWord string) {

	vk := api.Init(token)
	lp, err := longpoll.Init(vk, 2)
	if err != nil {
		log.Fatal(err)
	}

	regexp1 := regexp.MustCompile(fmt.Sprintf("^%v(-)?([0-9])?", strings.ToLower(triggerWord)))

	lp.EventNew(4, func(event []interface{}) error {
		var message Message
		message.Event = 4
		message.ID = int(event[1].(float64))
		message.Flags = int(event[2].(float64))
		message.PeerID = int(event[3].(float64))
		message.Timestamp = int(event[4].(float64))
		message.Text = strings.ToLower(event[5].(string))

		if (message.Flags & 1 << 1) == 0 {
			return nil
		}

		result := regexp1.FindStringSubmatch(message.Text)

		var (
			toDeleteReplace bool
			count           int
		)

		if result == nil {
			return nil
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
				if v != message.ID {
					_, err := vk.MessagesEdit(api.Params{"peer_id": message.PeerID, "message_id": v, "message": "ᅠ"})
					log.Println(err)
					switch errors.GetType(err) {
					case errors.Captcha:
						break
					}
				}
			}
			messages = append(messages, message.ID)

			_, err = vk.MessagesDelete(api.Params{"message_ids": ToArray(messages), "delete_for_all": 1})
			if err != nil {
				_, _ = vk.MessagesDelete(api.Params{"message_ids": ToArray(messages), "delete_for_all": 1})
			}
		} else {
			// Удаление сообщений с помощью execute
			DeleteExec(vk, count+1, message.PeerID)
		}

		return nil
	})

	// Запуск и автоподнятие
	for {
		_ = lp.Run()
		time.Sleep(time.Second * 10)
	}

}
