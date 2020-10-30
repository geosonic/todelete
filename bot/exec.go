/*
 * Copyright (c) 2020 GeoSonic. All rights reserved.
 */

package bot

import (
	"fmt"
	"log"

	"github.com/SevereCloud/vksdk/v2/api"
)

/*

Execute функции

*/

// За рефакторинг execute кода спасибо https://vk.com/notqb

// language=AS
const code = `
// Количество которое необходимо удалить
var delete_count = parseInt(Args.delete_count);

// peer_id диалога, в котором удаляем
var peer_id = parseInt(Args.peer_id); // int

// Переменная отсортированных сообщений
var message_ids = [];

// Получаем сообщения
var messages = API.messages.getHistory({peer_id: peer_id, count: 200}).items + API.messages.getHistory({peer_id: peer_id, count: 200, offset: 200}).items;
// Получаем ID аккаунта
var self_id = API.users.get()[0].id;
// Текущее время (подсказал https://vk.com/id370457723)
var time = API.utils.getServerTime();

while (messages.length > 0 && message_ids.length < delete_count) {
	// Переменная сообщения
	var message = messages.shift();

	// Находим свои сообщения, сравнивая свой ID
	// и которые возможно удалить для всех
	if (message.from_id == self_id && (time - message.date) < 86400) message_ids.push(message.id);
}

// Если Args.del существует, значит удаляем сообщения, иначе возвращаем их ID
if (Args.del) {
	// Если у нас есть сообщения, которые можно удалить,
	// тогда удаляем сообщения
	if (message_ids) {
		return API.messages.delete({message_ids: message_ids, delete_for_all: 1}).length;
	}
	// Иначе возвращаем 0
	return 0;
}

// Возвращаем сообщения, которые можно удалить
return message_ids;
`

var compressedCode string

// Сжатие execute кода для производительности
func init() {
	var err error
	compressedCode, err = CompressJS(code)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(compressedCode)
}

func DeleteExec(vk *api.VK, toDeleteCount, peerID int) (int, error) {
	var deleted int

	err := vk.ExecuteWithArgs(compressedCode, api.Params{"delete_count": toDeleteCount, "peer_id": peerID, "del": true}, &deleted)

	return deleted, err
}

func GetMessages(vk *api.VK, toDeleteCount, peerID int) ([]int, error) {
	var resp []int

	err := vk.ExecuteWithArgs(compressedCode, api.Params{"delete_count": toDeleteCount, "peer_id": peerID}, &resp)

	return resp, err
}
