## Запуск Docker контейнера с Golang

Для работы всех функций используется Go modules*  

Rest api для запуска и остановки контейнеров. Это демонстративная версия, просто показыват "как можно сделать"

- Подключение дирктории  
в данном примере подключается директория "rw"
Данный каталог специально содержит уже некоторые файлы, дтя того, чтобы можно было их открыть с помощью разаернутого контейнера.
- Компиляция и запуск программы
```sh
go get
go run .
```
- Запустим контейнер
```shell script
curl -X POST http://localhost:8180/kod/runContainer
```
Контейнер запускается, и консоль выводит "Container ID", он нужен для остановки. 
```shell script
curl -X POST http://localhost:8180/kod/stopContainer -d "your-container-id"
```
Зайти в развернутый контейнер kodexplorer можно по адресу http://localhost:8206  
Admin dir home:/var/www/html/data/User/admin/home/desktop/  
Файлы которые будут доступны всем пользователям:/data/Group/public/home/ProjectFiles  
При первом запуске попросит установить пароль администратора  
Уже существующие пользователи:  
demo|demo    
guest|guest    
Созданные при первом запуске файлы имеют специфичного владельца, что бы получить доступ можно установить соответствующие права
```shell script
cd kodResources/
sudo chmod -R 777 .
```
[![Demo docker api Xela golang](./dockerApiGolangXela.gif)](./dockerApiGolangXela.gif)

## Ручной запуск xaljer/kodexplorer
Запуск контейнера с командной строки  
```shell script
docker container run -d -p 8026:80 -v "/$(pwd)/rw:/var/www/html" xaljer/kodexplorer
```
Нельзя отредактировать папку rw, так как не являюсь владельцем, исправим  
```shell script
sudo chmod -R 777 .
```

Запустим контейнер с tmpfs
```shell script
curl -X POST http://localhost:8180/kod/runContainerTmpfs
```
C подключением разобрались

