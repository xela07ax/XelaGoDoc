
## Запускаем пример Gorilla
Источник: https://github.com/gorilla/websocket  

Связь между сервисами по протоколу Websocket или "Client and server example"

В этом примере показаны простые клиент и сервер.

Сервер повторяет отправленные ему сообщения. Клиент отправляет сообщение каждую секунду и распечатывает все полученные сообщения.

Чтобы запустить пример, запустите сервер:

    $ go run server.go

Далее запускаем клиент:

    $ go run client.go

Сервер включает в себя простой веб-клиент. Чтобы использовать клиент, откройте в браузере http://127.0.0.1:8080 и следуйте инструкциям на странице.