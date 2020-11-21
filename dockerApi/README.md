# Запуск Docker контейнера с Golang
## Описание
- Подключение локальной папки Volumes
- Открытие портов, а так же прокидывание на хост машину
- Управление контейнерами по id
- Управление контейнерами REST запросами через curl
- Демонстрация шикарного web файлового менеджера

Для работы всех функций используется Go modules*  

Rest api для запуска и остановки контейнеров. Это демонстративная версия, просто показывает "как можно сделать"

## Руководство запуска и использования
- Подключение директории  
в данном примере подключается директория "kodResources"
Данный каталог содержит используется для хранения профиля, и пользовательских данных.
- Компиляция и запуск программы
```sh
go get
go run .
```
- Запустим контейнер, а так-же прокинем внутренние порты контейнера на xост машину.
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
[![Demo docker api Xela golang](./_res/dockerApiGolangXela.gif)](https://youtu.be/FNgfSUm-P-4)

## Ручной запуск xaljer/kodexplorer
Запуск контейнера с командной строки  
```shell script
docker container run -d -p 8026:80 -v "/$(pwd)/rw:/var/www/html" xaljer/kodexplorer
```
Сам репозитарий: https://github.com/kalcaddle/KodExplorer  

## Теги:
docker api golang volume, create missing-type-opt: missing required option: "type", golang docker volume, docker api golang, go bind docker api, golang docker api file system, golang docker api, docker api map ports go, docker api golang list container, go modules локальный пакет, golang go-dockerclient, golang connect to docker
