/*
 * Copyright (c) 2020 GeoSonic. All rights reserved.
 */

package bot

import (
	"log"

	"github.com/SevereCloud/vksdk/api"
)

/*

Execute функции

*/

// За рефакторинг execute кода спасибо https://vk.com/notqb
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

while (messages.length > 0 && message_ids.length < delete_count) {
	// Переменная сообщения
	var message = messages.shift();

	// Находим свои сообщения
	if (message.from_id == self_id) {
		message_ids.push(message.id);
	}
}

// Если этот аргумент "type" равен 1, значит удаляем сообщения, иначе возвращаем их ID
if (Args.type == 1) {
	return API.messages.delete({message_ids: message_ids, delete_for_all: 1});
} else {
	return message_ids;
}
`

var compressedCode string

func init() {
	var err error
	compressedCode, err = CompressJS(code)
	if err != nil {
		log.Fatalln(err)
	}
}

func DeleteExec(vk *api.VK, toDeleteCount, peerID int) error {
	err := vk.ExecuteWithArgs(compressedCode, api.Params{"delete_count": toDeleteCount, "peer_id": peerID, "type": 1}, nil)

	return err
}

func GetMessages(vk *api.VK, toDeleteCount, peerID int) ([]int, error) {
	var resp []int

	err := vk.ExecuteWithArgs(compressedCode, api.Params{"delete_count": toDeleteCount, "peer_id": peerID, "type": 0}, &resp)

	return resp, err
}
