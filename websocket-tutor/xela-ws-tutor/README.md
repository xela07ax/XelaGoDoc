
## Коммуникация сервисов по Websocket на Golang
Видео: https://youtu.be/OiK8pKMqpWA

Два файла, это как две отдельных программы. В каждой из нах поднят свой HTTP сервис.
Программе client.go будем отправлять сообщение с помощью POST запроса, а server.go будет раболтать только по ws протоколу, он примет сообщение и отправит его обратно. На этом все завершится.

- server.go  
Этот файл мы представляем как некий "сервис", который должен где то крутиться, а мы к нему подключаемся.

- client.go  
Эту часть можно убрать куда нибудь в аркестратор, через него и будем отправлять команды.

Команда:
```sh
curl -H "Content-Type: application/json" -X POST http://localhost:7456/sh -d "Hello World"
```

## Теги
Вэбсокеты для Го самая эффективная связка, Websocker go example, ws golang start, ws go and go, общение Го сервисов по протоколу Websocket, go websocket client, websocket server

