# To Delete Bot
Скрипт на Golang, который позволяет удалять свои сообщения в ВК используя триггер-слово.

Функционал взят [отсюда](https://github.com/P2LOVE/VK-UserSide-Bot).
Используется библиотека [VK SDK](https://github.com/SevereCloud/vksdk) от [SevereCloud](https://github.com/SevereCloud).

![](example.gif)

# Установка

Для начала необходимо установить компилятор Golang, а также скачать библиотеку [VK SDK](https://github.com/SevereCloud/vksdk).

# Настройка для одного аккаунта

Необходимо, чтобы у токена был [доступ](https://vk.com/dev/messages_api) к методу messages. Я использовал токен от Kate Mobile, получить его можно тут https://vkhost.github.io.

main.go
```go
accounts := map[string]string{
	"токен": "триггер-слово",
}
```

# Несколько аккаунтов

main.go
```go
accounts := map[string]string{
	"токен": "триггер-слово",
	"токен": "триггер-слово",
	"токен": "триггер-слово", // и так далее
}
```

# Запуск

Сначала проинициализируем модули (Golang версии минимум 1.14).

```shell
go mod init deleter
```

После установки, настройки, скрипт готов к работе.

```shell
go run main.go
```

Также вы можете скомпилировать заранее и легко его запускать (можно будет запускать там где не будет компилятора).
В интернете можно найти гайды о том, как компилировать скрипты для другой платформы.

```shell
go build main.go
```

# Обратная связь

Все баги и предложения в Issues, Pull Requests также буду смотреть.
