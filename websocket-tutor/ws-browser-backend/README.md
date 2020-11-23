## Подключение к вебсокетам с браузера и с go приложения

srvMain.go - Запускает HTTP сервер и вэюсокеты, имеет несколько функциональных роутов:
- __"/"__ - по этому адресу можно зайти с браузера, подготовится страница с помощбю простейшего щаблонизатора, с этой страницы можно выполнить подключение по вэюсокетам на это-же приложение.
- 
```sh
go run srvMain.go
```
```sh
✗ go run tstBack.go
✗ curl -X POST http://localhost:3450/conn -d "localhost:3449"
✗ curl -X POST http://localhost:3450/send -d "Hellou"
```
```sh
✗ go run tstBack.go -addr=localhost:3456
✗ curl -X POST http://localhost:3456/conn -d "localhost:3449"
✗ curl -X POST http://localhost:3456/send -d "Hellou"
```
