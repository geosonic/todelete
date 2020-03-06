package bot

import (
	"fmt"
	"github.com/SevereCloud/vksdk/api"
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
// Получаем ID аккаунта
var myID = API.users.get()[0].id;
		
// Количество элементов для цикла
var count = resp.items.length;
// Счётчик для цикла
var counter = 0;
		
while (counter < count) {
// Переменная сообщения
var message = resp.items[counter];
		
if (message.from_id == myID && toDelete.length < toDeleteCount) {
toDelete.push(message.id);
}
		
// Итерация
counter = counter + 1;
}
`

func DeleteExec(vk *api.VK, toDeleteCount, peerID int) {
	code :=
		fmt.Sprintf(code+`// Возвращаем результат
		return API.messages.delete({message_ids: toDelete, delete_for_all: 1});`, toDeleteCount, peerID)

	_ = vk.Execute(code, nil)
}

func GetMessages(vk *api.VK, toDeleteCount, peerID int) []int {
	code := fmt.Sprintf(code+`// Возвращаем результат
		return {messages: toDelete};`, toDeleteCount, peerID)

	var resp struct {
		Messages []int `json:"messages"`
	}

	_ = vk.Execute(code, &resp)

	return resp.Messages

}
