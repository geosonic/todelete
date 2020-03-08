package bot

import (
	"fmt"
	"github.com/SevereCloud/vksdk/api"
	"strconv"
)

/*

Execute функции

*/

const code = `
// Количество которое необходимо удалить
// Эти переменные устанавливаются скриптом!!!
var toDeleteCount = %v; // int
var peer_id = %v; // int
		
// Переменная отсортированных сообщений
var toDelete = [];
		
// Получаем сообщения
var resp = API.messages.getHistory({peer_id: peer_id, count: 200});
resp.items = resp.items + API.messages.getHistory({peer_id: peer_id, count: 200, offset: 200}).items;
// Получаем ID аккаунта
var myID = API.users.get()[0].id;
		
// Количество элементов для цикла
var count = resp.items.length;
// Счётчик для цикла
var counter = 0;

while (counter < count) {
// Переменная сообщения
var message = resp.items[counter];

// Находим свои сообщения
if (message.from_id == myID && toDelete.length < toDeleteCount) {
toDelete.push(message.id);
}
		
// Итерация
counter = counter + 1;
}
`

func ToArray(slice []int) string {
	var s string

	for i := 0; i < len(slice); i++ {
		if i > 0 {
			s += ", "
		}
		s += strconv.Itoa(slice[i])
	}
	return s

}

func DeleteExec(vk *api.VK, toDeleteCount, peerID int) {
	if toDeleteCount > 99999999999 {
		toDeleteCount = 99999999999
	}
	code :=
		fmt.Sprintf(code+`// Возвращаем результат удаления сообщений
		return API.messages.delete({message_ids: toDelete, delete_for_all: 1});`, toDeleteCount, peerID)

	_ = vk.Execute(code, nil)
}

func GetMessages(vk *api.VK, toDeleteCount, peerID int) []int {
	if toDeleteCount > 99999999999 {
		toDeleteCount = 99999999999
	}
	code := fmt.Sprintf(code+`// Возвращаем найденные сообщения
		return {messages: toDelete};`, toDeleteCount, peerID)

	var resp struct {
		Messages []int `json:"messages"`
	}

	_ = vk.Execute(code, &resp)

	return resp.Messages
}
